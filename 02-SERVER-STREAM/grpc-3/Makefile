.PHONY: all client server clean

all: client server

gen: 
	protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb

clean: 
	rm -rf pb/.*.go

serve:
	cd server && go run .

client:
	cd client && go run .

help:
	@echo "gen: generate pb"
	@echo "clean: clean pb"
	@echo "serve: run server"
	@echo "client: run client"