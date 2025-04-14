
```markdown
# 🚀 Go gRPC Best Practices Guide

## 📋 Table of Contents
1. [Project Structure](#project-structure)
2. [Protocol Buffer Design](#protocol-buffer-design)
3. [Error Handling](#error-handling)
4. [Context Management & Deadlines](#context-management--deadlines)
5. [Interceptors](#interceptors)
6. [Service Discovery](#service-discovery)
7. [Validation](#validation)
8. [Testing](#testing)
9. [Middleware](#middleware)
10. [Streaming Best Practices](#streaming-best-practices)

## 📁 Project Structure

### Recommended Layout
```bash
/project-root
├── api/
│   └── proto/
│       ├── v1/
│       │   ├── user.proto
│       │   └── order.proto
│       └── v2/
├── internal/
│   ├── server/
│   │   └── grpc/
│   ├── service/
│   ├── repository/
│   └── domain/
└── pkg/
    ├── errors/
    ├── validation/
    └── middleware/
```

### Proto File Organization
```protobuf
// api/proto/v1/user.proto
syntax = "proto3";

package myapp.user.v1;
option go_package = "github.com/org/project/api/v1;userv1";

import "google/api/annotations.proto";
import "google/protobuf/timestamp.proto";
import "validate/validate.proto";

service UserService {
    rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {
        option (google.api.http) = {
            post: "/v1/users"
            body: "*"
        };
    }
}

message User {
    string id = 1 [(validate.rules).string.uuid = true];
    string email = 2 [(validate.rules).string.email = true];
    string full_name = 3 [(validate.rules).string = {
        min_len: 1
        max_len: 100
    }];
    repeated string roles = 4;
    google.protobuf.Timestamp created_at = 5;
}
```

## 🎯 Error Handling

### Custom Error Types
```go
// pkg/errors/errors.go
type ErrorType int

const (
    Unknown ErrorType = iota
    ValidationError
    NotFoundError
    AlreadyExistsError
)

type Error struct {
    Type    ErrorType
    Message string
    Code    codes.Code
    Err     error
}

func (e *Error) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func (e *Error) GRPCStatus() *status.Status {
    return status.New(e.Code, e.Message)
}

// Error constructor functions
func NewValidationError(msg string, err error) error {
    return &Error{
        Type:    ValidationError,
        Message: msg,
        Code:    codes.InvalidArgument,
        Err:     err,
    }
}
```

### Error Handling in Services
```go
// internal/service/user.go
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    if err := req.Validate(); err != nil {
        return nil, errors.NewValidationError("invalid request", err)
    }

    user, err := s.repo.CreateUser(ctx, req.GetUser())
    if err != nil {
        switch {
        case errors.Is(err, repository.ErrDuplicate):
            return nil, status.Error(codes.AlreadyExists, "user already exists")
        default:
            s.logger.Error("failed to create user", zap.Error(err))
            return nil, status.Error(codes.Internal, "internal error")
        }
    }

    return &pb.CreateUserResponse{
        User: user,
    }, nil
}
```

## ⏰ Context Management & Deadlines

### Context Usage
```go
// internal/service/user.go
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    // Add timeout to context
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    // Check context before expensive operations
    if err := ctx.Err(); err != nil {
        return nil, status.Error(codes.DeadlineExceeded, "request timeout")
    }

    user, err := s.repo.GetUser(ctx, req.GetId())
    if err != nil {
        return nil, err
    }

    return user, nil
}
```

### Deadline Propagation
```go
// internal/service/order.go
func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
    // Get deadline from parent context
    deadline, ok := ctx.Deadline()
    if !ok {
        // Set default deadline if none exists
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
        defer cancel()
    } else {
        // Subtract time for local processing
        timeout := time.Until(deadline) - 100*time.Millisecond
        if timeout <= 0 {
            return nil, status.Error(codes.DeadlineExceeded, "insufficient time to process")
        }
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, timeout)
        defer cancel()
    }

    // Proceed with operation
    return s.processOrder(ctx, req)
}
```

## 🔄 Interceptors

### Unary Interceptors
```go
// pkg/middleware/interceptors.go
func LoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler) (interface{}, error) {
        
        start := time.Now()
        resp, err := handler(ctx, req)
        duration := time.Since(start)

        // Extract trace ID if using tracing
        var traceID string
        if span := trace.SpanFromContext(ctx); span != nil {
            traceID = span.SpanContext().TraceID().String()
        }

        logger.Info("handled request",
            zap.String("method", info.FullMethod),
            zap.Duration("duration", duration),
            zap.String("trace_id", traceID),
            zap.Error(err),
        )

        return resp, err
    }
}

func ValidationInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler) (interface{}, error) {
        
        if v, ok := req.(interface{ Validate() error }); ok {
            if err := v.Validate(); err != nil {
                return nil, status.Error(codes.InvalidArgument, err.Error())
            }
        }

        return handler(ctx, req)
    }
}
```

### Stream Interceptors
```go
func StreamLoggingInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
    return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo,
        handler grpc.StreamHandler) error {
        
        start := time.Now()
        err := handler(srv, ss)
        duration := time.Since(start)

        logger.Info("handled stream",
            zap.String("method", info.FullMethod),
            zap.Duration("duration", duration),
            zap.Error(err),
        )

        return err
    }
}
```

## 🔍 Service Discovery

### Service Registry Interface
```go
// pkg/registry/registry.go
type ServiceInfo struct {
    Name     string
    Version  string
    Address  string
    Port     int
    Methods  []string
}

type ServiceRegistry interface {
    Register(service *ServiceInfo) error
    Deregister(serviceName string) error
    GetService(name string) ([]*ServiceInfo, error)
    Watch(name string) (<-chan []*ServiceInfo, error)
}
```

### Service Discovery Implementation
```go
// internal/server/grpc/server.go
type Server struct {
    registry ServiceRegistry
    services map[string]*ServiceInfo
}

func (s *Server) RegisterService(info *ServiceInfo) error {
    if err := s.registry.Register(info); err != nil {
        return fmt.Errorf("failed to register service: %w", err)
    }

    s.services[info.Name] = info
    return nil
}

// Service client with discovery
func NewServiceClient(ctx context.Context, registry ServiceRegistry, serviceName string) (*grpc.ClientConn, error) {
    services, err := registry.GetService(serviceName)
    if err != nil {
        return nil, fmt.Errorf("failed to get service: %w", err)
    }

    // Use first available service instance
    service := services[0]
    conn, err := grpc.DialContext(ctx,
        fmt.Sprintf("%s:%d", service.Address, service.Port),
        grpc.WithInsecure(),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to dial service: %w", err)
    }

    return conn, nil
}
```

## ✅ Validation

### Request Validation
```go
// pkg/validation/validator.go
type Validator interface {
    Validate() error
}

func ValidateRequest(req interface{}) error {
    if v, ok := req.(Validator); ok {
        if err := v.Validate(); err != nil {
            return errors.NewValidationError("invalid request", err)
        }
    }
    return nil
}

// Usage in service
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
    if err := validation.ValidateRequest(req); err != nil {
        return nil, err
    }
    // Process request
}
```

## 🧪 Testing

### Unit Testing Services
```go
// internal/service/user_test.go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        req     *pb.CreateUserRequest
        mock    func(*mockRepo)
        want    *pb.User
        wantErr codes.Code
    }{
        {
            name: "successful creation",
            req: &pb.CreateUserRequest{
                Email: "test@example.com",
                Name:  "Test User",
            },
            mock: func(m *mockRepo) {
                m.EXPECT().
                    CreateUser(gomock.Any(), gomock.Any()).
                    Return(&pb.User{
                        Id:    "123",
                        Email: "test@example.com",
                        Name:  "Test User",
                    }, nil)
            },
            want: &pb.User{
                Id:    "123",
                Email: "test@example.com",
                Name:  "Test User",
            },
            wantErr: codes.OK,
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockRepo := NewMockRepository(ctrl)
            tt.mock(mockRepo)

            svc := NewUserService(mockRepo)
            got, err := svc.CreateUser(context.Background(), tt.req)

            if err != nil {
                if st, ok := status.FromError(err); ok {
                    assert.Equal(t, tt.wantErr, st.Code())
                } else {
                    t.Errorf("error is not a status error: %v", err)
                }
                return
            }

            assert.Equal(t, tt.want, got)
        })
    }
}
```

### Integration Testing
```go
// internal/service/integration_test.go
func TestIntegration_UserService(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    ctx := context.Background()

    // Setup test server
    srv, err := setupTestServer(t)
    require.NoError(t, err)
    defer srv.Stop()

    // Create client
    conn, err := grpc.DialContext(ctx, srv.Addr,
        grpc.WithInsecure(),
        grpc.WithBlock(),
    )
    require.NoError(t, err)
    defer conn.Close()

    client := pb.NewUserServiceClient(conn)

    // Run tests
    t.Run("create and get user", func(t *testing.T) {
        // Create user
        createResp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
            Email: "test@example.com",
            Name:  "Test User",
        })
        require.NoError(t, err)

        // Get user
        getResp, err := client.GetUser(ctx, &pb.GetUserRequest{
            Id: createResp.User.Id,
        })
        require.NoError(t, err)

        assert.Equal(t, createResp.User.Email, getResp.Email)
        assert.Equal(t, createResp.User.Name, getResp.Name)
    })
}
```

## 📡 Streaming Best Practices

### Client Streaming
```go
// internal/service/user.go
func (s *UserService) UploadUserDocuments(stream pb.UserService_UploadUserDocumentsServer) error {
    ctx := stream.Context()
    var documents []*pb.Document

    for {
        // Check context before receiving
        if err := ctx.Err(); err != nil {
            return status.Error(codes.Canceled, "stream canceled")
        }

        req, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return status.Errorf(codes.Internal, "failed to receive: %v", err)
        }

        documents = append(documents, req.GetDocument())
    }

    // Process documents
    result, err := s.processDocuments(ctx, documents)
    if err != nil {
        return err
    }

    return stream.SendAndClose(result)
}
```

### Server Streaming
```go
func (s *UserService) WatchUserActivity(req *pb.WatchUserRequest, stream pb.UserService_WatchUserActivityServer) error {
    ctx := stream.Context()
    activities := make(chan *pb.UserActivity, 100)

    // Start watching activities
    go s.watchActivities(ctx, req.GetUserId(), activities)

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case activity, ok := <-activities:
            if !ok {
                return nil
            }
            if err := stream.Send(activity); err != nil {
                return status.Errorf(codes.Internal, "failed to send: %v", err)
            }
        }
    }
}
```

### Bi-directional Streaming
```go
func (s *UserService) ChatSupport(stream pb.UserService_ChatSupportServer) error {
    ctx := stream.Context()
    
    // Create message channels
    incoming := make(chan *pb.ChatMessage, 100)
    outgoing := make(chan *pb.ChatMessage, 100)

    // Start processing incoming messages
    go func() {
        defer close(incoming)
        for {
            msg, err := stream.Recv()
            if err == io.EOF {
                return
            }
            if err != nil {
                s.logger.Error("failed to receive message", zap.Error(err))
                return
            }
            incoming <- msg
        }
    }()

    // Start chat processing
    go s.processChatMessages(ctx, incoming, outgoing)

    // Send outgoing messages
    for msg := range outgoing {
        if err := stream.Send(msg); err != nil {
            return status.Errorf(codes.Internal, "failed to send: %v", err)
        }
    }

    return nil
}
```

## 🔒 Middleware

### Authentication Middleware
```go
// pkg/middleware/auth.go
func AuthInterceptor(authClient auth.Client) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler) (interface{}, error) {
        
        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return nil, status.Error(codes.Unauthenticated, "missing metadata")
        }

        token := extractToken(md)
        if token == "" {
            return nil, status.Error(codes.Unauthenticated, "missing token")
        }

        claims, err := authClient.ValidateToken(ctx, token)
        if err != nil {
            return nil, status.Error(codes.Unauthenticated, "invalid token")
        }

        // Add claims to context
        newCtx := context.WithValue(ctx, "claims", claims)
        return handler(newCtx, req)
    }
}
```

### Rate Limiting Middleware
```go
// pkg/middleware/ratelimit.go
func RateLimitInterceptor(limit rate.Limit, burst int) grpc.UnaryServerInterceptor {
    limiter := rate.NewLimiter(limit, burst)
    
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler) (interface{}, error) {
        
        if err := limiter.Wait(ctx); err != nil {
            return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
        }

        return handler(ctx, req)
    }
}
```

### Panic Recovery Middleware
```go
// pkg/middleware/recovery.go
func RecoveryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler) (resp interface{}, err error) {
        
        defer func() {
            if r := recover(); r != nil {
                logger.Error("panic recovered",
                    zap.Any("panic", r),
                    zap.String("stack", string(debug.Stack())),
                )
                err = status.Error(codes.Internal, "internal error")
            }
        }()

        return handler(ctx, req)
    }
}
```