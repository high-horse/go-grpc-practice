```markdown
# 📘 gRPC Go Best Practices Guide

[![Go Reference](https://pkg.go.dev/badge/google.golang.org/grpc.svg)](https://pkg.go.dev/google.golang.org/grpc)

## 📋 Table of Contents
- [Project Structure](#project-structure)
- [Protobuf Guidelines](#protobuf-guidelines)
- [Service Implementation](#service-implementation)
- [Error Handling](#error-handling)
- [Security](#security)
- [Testing](#testing)
- [Performance](#performance)

## 🏗️ Project Structure
```bash
.
├── api/
│   └── proto/              # Protocol buffer definitions
│       ├── v1/
│       └── v2/
├── cmd/                    # Main applications
├── internal/              # Private application code
├── pkg/                   # Public library code
└── tests/                # Integration tests
```

## 📝 Protobuf Guidelines

### Versioning
```protobuf
// File: api/proto/v1/user.proto
syntax = "proto3";

package myapp.v1;
option go_package = "myapp/api/v1;v1";
```

### Service Definition
```protobuf
service UserService {
  // Use meaningful names and clear documentation
  rpc GetUser(GetUserRequest) returns (GetUserResponse) {
    option (google.api.http) = {
      get: "/v1/users/{user_id}"
    };
  }
  
  // Use streaming where appropriate
  rpc WatchUserStatus(WatchUserRequest) returns (stream UserStatus) {}
}
```

## ⚙️ Service Implementation

### Server Setup
```go
// Best practice server initialization
func NewGRPCServer(opts ...Option) (*grpc.Server, error) {
    // Default interceptors
    interceptors := []grpc.UnaryServerInterceptor{
        grpc_recovery.UnaryServerInterceptor(),
        grpc_ctxtags.UnaryServerInterceptor(),
        grpc_prometheus.UnaryServerInterceptor,
    }
    
    // Server options
    serverOpts := []grpc.ServerOption{
        grpc.ChainUnaryInterceptor(interceptors...),
    }
    
    return grpc.NewServer(serverOpts...), nil
}
```

### Interceptors
```go
func LoggingInterceptor(ctx context.Context, req interface{}, 
    info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    
    start := time.Now()
    resp, err := handler(ctx, req)
    
    logger.Info("request processed",
        "method", info.FullMethod,
        "duration", time.Since(start),
        "error", err)
        
    return resp, err
}
```

## 🚨 Error Handling

### Error Types
```go
// Custom error types
type NotFoundError struct {
    Resource string
    ID       string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}

// Error handling in services
func (s *service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    user, err := s.repo.GetUser(ctx, req.GetId())
    if err != nil {
        switch {
        case errors.Is(err, &NotFoundError{}):
            return nil, status.Error(codes.NotFound, err.Error())
        default:
            return nil, status.Error(codes.Internal, "internal error")
        }
    }
    return user, nil
}
```

## 🔒 Security

### TLS Configuration
```go
func setupTLSCredentials() (credentials.TransportCredentials, error) {
    cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
    if err != nil {
        return nil, err
    }
    
    config := &tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientAuth:   tls.RequireAndVerifyClientCert,
    }
    
    return credentials.NewTLS(config), nil
}
```

### Authentication
```go
func AuthInterceptor(ctx context.Context, req interface{}, 
    info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    
    token, err := extractToken(ctx)
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, "invalid token")
    }
    
    claims, err := validateToken(token)
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, "invalid token")
    }
    
    newCtx := context.WithValue(ctx, "user", claims)
    return handler(newCtx, req)
}
```

## 🧪 Testing

### Unit Tests
```go
func TestUserService_GetUser(t *testing.T) {
    tests := []struct {
        name    string
        userID  string
        mock    func(*mockRepo)
        want    *pb.User
        wantErr codes.Code
    }{
        {
            name:   "successful retrieval",
            userID: "123",
            mock: func(m *mockRepo) {
                m.EXPECT().GetUser(gomock.Any(), "123").Return(&pb.User{
                    Id:   "123",
                    Name: "Test User",
                }, nil)
            },
            want: &pb.User{
                Id:   "123",
                Name: "Test User",
            },
            wantErr: codes.OK,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests
```go
func TestIntegration_UserService(t *testing.T) {
    // Start test server
    srv, listener := setupTestServer(t)
    defer srv.Stop()
    
    // Create client
    conn, err := grpc.Dial(listener.Addr().String(), 
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        t.Fatalf("failed to dial: %v", err)
    }
    defer conn.Close()
    
    client := pb.NewUserServiceClient(conn)
    
    // Run tests
    t.Run("get user", func(t *testing.T) {
        // Test implementation
    })
}
```

## 🚀 Performance

### Connection Management
```go
func NewClientConn(ctx context.Context, target string) (*grpc.ClientConn, error) {
    opts := []grpc.DialOption{
        grpc.WithDefaultServiceConfig(`{
            "loadBalancingPolicy": "round_robin",
            "healthCheckConfig": {
                "serviceName": ""
            }
        }`),
        grpc.WithDefaultCallOptions(
            grpc.MaxCallRecvMsgSize(1024*1024*4), // 4MB
            grpc.MaxCallSendMsgSize(1024*1024*4),
        ),
    }
    
    return grpc.DialContext(ctx, target, opts...)
}
```

### Streaming Best Practices
```go
func (s *server) StreamData(req *pb.StreamRequest, stream pb.Service_StreamDataServer) error {
    // Set buffer size for channel
    buffer := make(chan *pb.Data, 100)
    
    // Handle context cancellation
    ctx := stream.Context()
    
    go func() {
        defer close(buffer)
        // Produce data
    }()
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case data, ok := <-buffer:
            if !ok {
                return nil
            }
            if err := stream.Send(data); err != nil {
                return err
            }
        }
    }
}
```

## 📈 Monitoring

### Prometheus Metrics
```go
func setupMetrics(srv *grpc.Server) {
    grpc_prometheus.Register(srv)
    http.Handle("/metrics", promhttp.Handler())
}
```

---

## 📚 Additional Resources
- [gRPC Go Documentation](https://grpc.io/docs/languages/go/)
- [Protocol Buffers Documentation](https://developers.google.com/protocol-buffers)
- [Go Style Guide](https://google.github.io/styleguide/go/)

```