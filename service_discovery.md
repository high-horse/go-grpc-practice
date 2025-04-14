
```markdown
# 🌐 gRPC Service Discovery & Inter-Service Communication

## 📋 Table of Contents
- [Service Architecture](#service-architecture)
- [Service Discovery Implementation](#service-discovery-implementation)
- [Service Communication](#service-communication)
- [Load Balancing](#load-balancing)
- [Health Checking](#health-checking)

## 🏗️ Service Architecture

### Service Registry Structure
```go
// services/registry/types.go
type ServiceInfo struct {
    Name        string
    Version     string
    ID          string
    Address     string
    Port        int
    Metadata    map[string]string
    Endpoints   []string
    Dependencies []string
}

type ServiceRegistry interface {
    Register(service *ServiceInfo) error
    Deregister(serviceID string) error
    GetService(name string) ([]*ServiceInfo, error)
    Watch(name string) (<-chan []*ServiceInfo, error)
}
```

### Example Microservices Structure
```bash
/microservices
├── /user-service
├── /order-service
├── /payment-service
├── /notification-service
└── /shared
    ├── /proto
    ├── /registry
    └── /middleware
```

## 🔄 Service Discovery Implementation

### Service Registration
```go
// services/shared/registry/consul.go
type ConsulRegistry struct {
    client *consul.Client
    logger *zap.Logger
}

func (r *ConsulRegistry) Register(service *ServiceInfo) error {
    registration := &consul.AgentServiceRegistration{
        ID:      service.ID,
        Name:    service.Name,
        Tags:    []string{service.Version},
        Port:    service.Port,
        Address: service.Address,
        Meta: map[string]string{
            "endpoints":    strings.Join(service.Endpoints, ","),
            "dependencies": strings.Join(service.Dependencies, ","),
        },
        Check: &consul.AgentServiceCheck{
            GRPC:                           fmt.Sprintf("%s:%d", service.Address, service.Port),
            GRPCUseTLS:                     false,
            Interval:                       "10s",
            Timeout:                        "5s",
            DeregisterCriticalServiceAfter: "30s",
        },
    }

    return r.client.Agent().ServiceRegister(registration)
}
```

### Service Implementation Example
```go
// services/user-service/main.go
type UserService struct {
    registry    registry.ServiceRegistry
    orderClient order.OrderServiceClient
    paymentClient payment.PaymentServiceClient
    logger     *zap.Logger
    pb.UnimplementedUserServiceServer
}

func NewUserService(registry registry.ServiceRegistry) *UserService {
    return &UserService{
        registry: registry,
        logger:   zap.L(),
    }
}

func (s *UserService) Start() error {
    // Register service
    serviceInfo := &registry.ServiceInfo{
        Name:    "user-service",
        Version: "v1",
        ID:      uuid.New().String(),
        Address: "localhost",
        Port:    50051,
        Endpoints: []string{
            "CreateUser",
            "GetUser",
            "UpdateUser",
        },
        Dependencies: []string{
            "order-service",
            "payment-service",
        },
    }

    if err := s.registry.Register(serviceInfo); err != nil {
        return fmt.Errorf("failed to register service: %w", err)
    }

    // Initialize client connections
    if err := s.initClientConnections(); err != nil {
        return fmt.Errorf("failed to init clients: %w", err)
    }

    return nil
}
```

## 🔌 Service Communication

### Client Connection Manager
```go
// services/shared/client/manager.go
type ClientManager struct {
    registry registry.ServiceRegistry
    conns    map[string]*grpc.ClientConn
    mu       sync.RWMutex
}

func (m *ClientManager) GetConnection(serviceName string) (*grpc.ClientConn, error) {
    m.mu.RLock()
    if conn, ok := m.conns[serviceName]; ok {
        m.mu.RUnlock()
        return conn, nil
    }
    m.mu.RUnlock()

    m.mu.Lock()
    defer m.mu.Unlock()

    // Double-check after acquiring write lock
    if conn, ok := m.conns[serviceName]; ok {
        return conn, nil
    }

    // Get service instances
    instances, err := m.registry.GetService(serviceName)
    if err != nil {
        return nil, fmt.Errorf("failed to get service %s: %w", serviceName, err)
    }

    // Create connection with service discovery resolver
    conn, err := grpc.Dial(
        fmt.Sprintf("discovery:///%s", serviceName),
        grpc.WithDefaultServiceConfig(`{
            "loadBalancingPolicy": "round_robin",
            "healthCheckConfig": {
                "serviceName": ""
            }
        }`),
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create connection: %w", err)
    }

    m.conns[serviceName] = conn
    return conn, nil
}
```

### Inter-Service Communication Example
```go
// services/user-service/handlers.go
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    // Create user
    user, err := s.repository.CreateUser(ctx, req)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
    }

    // Create initial order using order service
    orderConn, err := s.clientManager.GetConnection("order-service")
    if err != nil {
        s.logger.Error("failed to get order service connection", zap.Error(err))
        return nil, status.Error(codes.Internal, "internal error")
    }
    
    orderClient := order.NewOrderServiceClient(orderConn)
    orderResp, err := orderClient.CreateOrder(ctx, &order.CreateOrderRequest{
        UserId: user.Id,
    })
    if err != nil {
        s.logger.Error("failed to create initial order", zap.Error(err))
        // Handle error but don't fail the user creation
    }

    return &pb.CreateUserResponse{
        User: user,
        InitialOrderId: orderResp.GetOrderId(),
    }, nil
}
```

## ⚖️ Load Balancing

### Custom Load Balancer
```go
// services/shared/loadbalancer/balancer.go
type ServiceLoadBalancer struct {
    services  []*registry.ServiceInfo
    index    uint64
}

func (lb *ServiceLoadBalancer) Next() *registry.ServiceInfo {
    if len(lb.services) == 0 {
        return nil
    }
    
    idx := atomic.AddUint64(&lb.index, 1)
    return lb.services[idx%uint64(len(lb.services))]
}

func (lb *ServiceLoadBalancer) UpdateServices(services []*registry.ServiceInfo) {
    lb.services = services
}
```

### Service Discovery Resolver
```go
// services/shared/resolver/discovery.go
type discoveryResolver struct {
    target     resolver.Target
    cc         resolver.ClientConn
    registry   registry.ServiceRegistry
    balancer   *ServiceLoadBalancer
    updateChan chan struct{}
}

func (r *discoveryResolver) watch() {
    watchChan, err := r.registry.Watch(r.target.Endpoint)
    if err != nil {
        r.cc.ReportError(err)
        return
    }

    for {
        select {
        case services := <-watchChan:
            r.balancer.UpdateServices(services)
            r.updateAddresses()
        case <-r.updateChan:
            return
        }
    }
}

func (r *discoveryResolver) updateAddresses() {
    services, err := r.registry.GetService(r.target.Endpoint)
    if err != nil {
        r.cc.ReportError(err)
        return
    }

    addresses := make([]resolver.Address, len(services))
    for i, svc := range services {
        addresses[i] = resolver.Address{
            Addr:       fmt.Sprintf("%s:%d", svc.Address, svc.Port),
            ServerName: svc.Name,
            Attributes: attributes.New(
                "version", svc.Version,
                "metadata", svc.Metadata,
            ),
        }
    }

    r.cc.UpdateState(resolver.State{Addresses: addresses})
}
```

## 🏥 Health Checking

### Health Check Service
```go
// services/shared/health/service.go
type HealthService struct {
    registry registry.ServiceRegistry
    checker  health.HealthCheck
}

func (s *HealthService) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
    service, err := s.registry.GetService(req.Service)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "service %s not found", req.Service)
    }

    status := health.HealthCheckResponse_SERVING
    for _, svc := range service {
        if err := s.checker.Check(svc); err != nil {
            status = health.HealthCheckResponse_NOT_SERVING
            break
        }
    }

    return &health.HealthCheckResponse{
        Status: status,
    }, nil
}
```

## 📊 Monitoring & Tracing

### Service Metrics
```go
// services/shared/metrics/metrics.go
var (
    serviceRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "grpc_service_requests_total",
            Help: "Total number of gRPC requests by service and method",
        },
        []string{"service", "method"},
    )

    serviceLatency = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "grpc_service_latency_seconds",
            Help:    "Service method latency in seconds",
            Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2.5, 5},
        },
        []string{"service", "method"},
    )
)

func RecordMetrics(ctx context.Context, service, method string) func() {
    start := time.Now()
    serviceRequests.WithLabelValues(service, method).Inc()

    return func() {
        duration := time.Since(start).Seconds()
        serviceLatency.WithLabelValues(service, method).Observe(duration)
    }
}
```

### Distributed Tracing
```go
// services/shared/tracing/tracer.go
func InitTracer(serviceName string) (*trace.TracerProvider, error) {
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
    ))
    if err != nil {
        return nil, fmt.Errorf("failed to create jaeger exporter: %w", err)
    }

    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(serviceName),
        )),
    )

    otel.SetTracerProvider(tp)
    return tp, nil
}
```

## 🚀 Service Launch Example

```go
// services/user-service/main.go
func main() {
    // Initialize logger
    logger, _ := zap.NewProduction()
    defer logger.Sync()
    zap.ReplaceGlobals(logger)

    // Initialize tracer
    tp, err := tracing.InitTracer("user-service")
    if err != nil {
        logger.Fatal("failed to init tracer", zap.Error(err))
    }
    defer tp.Shutdown(context.Background())

    // Initialize registry
    reg, err := registry.NewConsulRegistry(&registry.Config{
        Address: "localhost:8500",
    })
    if err != nil {
        logger.Fatal("failed to create registry", zap.Error(err))
    }

    // Create service
    svc := NewUserService(reg)
    if err := svc.Start(); err != nil {
        logger.Fatal("failed to start service", zap.Error(err))
    }

    // Create gRPC server
    server := grpc.NewServer(
        grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
            grpc_prometheus.UnaryServerInterceptor,
            grpc_zap.UnaryServerInterceptor(logger),
            grpc_recovery.UnaryServerInterceptor(),
        )),
    )

    // Register service
    pb.RegisterUserServiceServer(server, svc)
    health.RegisterHealthServer(server, health.NewServer())
    grpc_prometheus.Register(server)

    // Start server
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        logger.Fatal("failed to listen", zap.Error(err))
    }

    logger.Info("starting server", zap.String("address", lis.Addr().String()))
    if err := server.Serve(lis); err != nil {
        logger.Fatal("failed to serve", zap.Error(err))
    }
}
```

## 🔒 Security Considerations

### Service-to-Service Authentication
```go
// services/shared/auth/middleware.go
func ServiceAuthInterceptor(authKey string) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler) (interface{}, error) {
        
        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return nil, status.Error(codes.Unauthenticated, "missing metadata")
        }

        values := md.Get("service-auth-key")
        if len(values) == 0 || values[0] != authKey {
            return nil, status.Error(codes.Unauthenticated, "invalid service auth key")
        }

        return handler(ctx, req)
    }
}
```

## 📝 Configuration Example

### Service Configuration
```yaml
# config/user-service.yaml
service:
  name: user-service
  version: v1.0.0
  port: 50051

registry:
  type: consul
  address: localhost:8500
  check_interval: 10s
  deregister_after: 30s

dependencies:
  - name: order-service
    version: v1
    required: true
  - name: payment-service
    version: v1
    required: false

security:
  service_auth_key: ${SERVICE_AUTH_KEY}
  tls_enabled: true
  cert_file: /etc/certs/service.crt
  key_file: /etc/certs/service.key

tracing:
  enabled: true
  jaeger_endpoint: http://jaeger:14268/api/traces
  sample_rate: 0.1

metrics:
  enabled: true
  prometheus_port: 9090
```