gen :
	protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb

clean :
	rm pb/*.go

serve :
	go run server/main.go

client : 	
	go run client/main.go