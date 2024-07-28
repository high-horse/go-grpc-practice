package main

import (
	"log"
	"net"

	proto "grpc-1/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const PORT = ":50051"


func main() {
	println("Server starting at ", PORT)
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal(err)
	}

	// create a new grpc server
	// serv := grpc.NewServer()
	serv := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryInterseptor),
		grpc.StreamInterceptor(StreamInterseptor),
	)

	// Register the service with the server
	proto.RegisterNewserviceServer(serv, &Server{})

	/*
	Without reflection, 
	you need the .proto files to understand what services and methods a gRPC server offers. 
	With reflection, 
	tools can query the server directly to get this information, 
	making development and debugging easier.
	*/

	// Register reflection service on gRPC server (optional, for debugging)
	reflection.Register(serv)

	if err := serv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}



