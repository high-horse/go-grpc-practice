.PHONY: clean server clinet

gen:
	protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb

clean:
	rm -rf pb/.*.go

server:
	cd server && go run .

client: 
	cd client && go run .