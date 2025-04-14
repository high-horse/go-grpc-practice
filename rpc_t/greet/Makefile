PROTO_DIR=protos
PROTO_FILES=$(shell find $(PROTO_DIR) -name "*.proto")

gen_all:
	@echo "Generating proto code..."
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(PROTO_DIR) \
		--go-grpc_out=$(PROTO_DIR) \
		$(PROTO_FILES)
	@echo "Proto code generation completed."

gen_greet:
	@echo "Generating proto code..."
	protoc   --proto_path=./protos   --go_out=.   --go-grpc_out=.   api/greet.proto
	@echo "Successfully generated code."
