.PHONY: gen clean server client sqlc-i goose-i goose-up goose-down sqlc help clean-sqlc

# Generate gRPC code from proto files
gen:
	protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb

# Clean generated gRPC code
clean:
	rm -rf pb/*.go

# Run the server
server:
	cd server && go run .

# Run the client
client:
	cd client-dummy && go run .

# Install sqlc and goose tools
sqlc-i:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest


# Run Goose migrations up
goose-up:
	cd store/schema && goose  postgres "postgres://postgres:root@localhost:5432/attempt2" up

# Run Goose migrations down
goose-down:
	cd store/schema && goose  postgres "postgres://postgres:root@localhost:5432/attempt2" down

# Generate code from SQL queries using sqlc
sqlc:
	sqlc generate

clean-sqlc:
	cd store && rm -rf database/*.go

# Display help information
help:
	@echo "Available targets:"
	@echo "  gen         - Generate gRPC code from proto files"
	@echo "  clean       - Clean generated gRPC code"
	@echo "  server      - Run the server"
	@echo "  client      - Run the client"
	@echo "  sqlc-i      - Install sqlc and goose tools"
	@echo "  goose-up    - Run Goose migrations up"
	@echo "  goose-down  - Run Goose migrations down"
	@echo "  sqlc        - Generate code from SQL queries using sqlc"
	@echo "  help        - Display this help information"
