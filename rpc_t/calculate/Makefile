gen :
	@echo "generating proto code ..."
	protoc   --proto_path=./protos   --go_out=.   --go-grpc_out=.   calculate/calculate.proto
	# protoc  --proto_path=./protos   --go_out=.   --go-grpc_out=.   calculate/calculate.proto
	@echo "successfully completed."