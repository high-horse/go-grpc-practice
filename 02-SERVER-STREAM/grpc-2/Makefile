.PHONY: gen client serve clean

gen :
	protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb

clean :
	rm pb/*.go

serve :
	go run $(wildcard server/*.go)

client : 
	@echo "Running client..."	
	go run  $(wildcard client/*.go)
