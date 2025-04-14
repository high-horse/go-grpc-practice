```bash
#!/bin/bash

# Colors for formatting
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Go & gRPC Best Practices Guide ===${NC}\n"

print_section() {
    echo -e "${BLUE}$1${NC}"
    echo -e "$2\n"
}

print_code_example() {
    echo -e "${RED}Example:${NC}"
    echo -e "$1\n"
}

print_section "Project Structure" "
├── cmd/                 # Main applications
├── internal/           # Private application code
├── pkg/                # Public library code
├── api/               # Protocol definition files
│   └── proto/         # Protocol buffers
├── configs/           # Configuration files
└── test/              # Additional external test apps/tools"

print_section "Go Code Organization" "
1. Use meaningful package names
2. Keep packages focused and cohesive
3. Avoid package cycles
4. Follow standard Go project layout"

print_code_example "
# Good package structure
/myproject
  /pkg
    /models
    /services
    /handlers
  /internal
    /database
    /auth"

print_section "Error Handling" "
1. Always check errors
2. Custom error types for specific cases
3. Use error wrapping
4. Return errors rather than panic"

print_code_example "
if err != nil {
    return fmt.Errorf(\"failed to process request: %w\", err)
}"

print_section "gRPC Service Definition" "
1. Use proto3 syntax
2. Clear service and message naming
3. Proper field numbering
4. Well-documented services"

print_code_example "
syntax = \"proto3\";

service UserService {
    rpc GetUser (GetUserRequest) returns (User) {}
    rpc ListUsers (ListUsersRequest) returns (stream User) {}
}

message User {
    string id = 1;
    string name = 2;
    string email = 3;
}"

print_section "gRPC Implementation" "
1. Use interceptors for cross-cutting concerns
2. Implement proper error handling
3. Use context for deadlines/cancellation
4. Include proper logging"

print_code_example "
func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    if err := validateRequest(req); err != nil {
        return nil, status.Error(codes.InvalidArgument, err.Error())
    }
    // Implementation
}"

print_section "Testing" "
1. Write unit tests for packages
2. Integration tests for services
3. Use test tables
4. Mock external dependencies"

print_code_example "
func TestUserService(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *User
        wantErr bool
    }{
        // test cases
    }
    // test implementation
}"

print_section "Performance" "
1. Use connection pooling
2. Implement proper timeouts
3. Use streaming where appropriate
4. Monitor and profile services"

print_section "Security" "
1. Use TLS for connections
2. Implement proper authentication
3. Validate all inputs
4. Use secure configuration management"

print_code_example "
creds := credentials.NewTLS(&tls.Config{})
server := grpc.NewServer(grpc.Creds(creds))"

echo -e "${GREEN}=== End of Best Practices Guide ===${NC}"
```

Save this as `go_grpc_best_practices.sh` and run:

```bash
chmod +x go_grpc_best_practices.sh
./go_grpc_best_practices.sh
```