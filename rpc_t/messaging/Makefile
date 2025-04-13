gen:
	@echo "generating code from proto file ..."
	protoc  --proto_path=./protos --go_out=.   --go-grpc_out=.  message/message.proto
	@echo "Proto code generation completed."