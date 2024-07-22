package main

import (
	proto "grpc-3/pb"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


const URL = "https://fakestoreapi.com/products"
const ADDRESS = ":50051"


func main() {	
	println("server starting at " + ADDRESS)

	listener, err := net.Listen("tcp", ADDRESS)
	if err != nil {
		log.Fatalf("Couldnot start server \n", err)
	}

	serv := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterseptor),
	)
	proto.RegisterProductServiceServer(serv, &Server{})

	reflection.Register(serv)

	if err := serv.Serve(listener); err != nil {
		log.Fatalf("Could not start server \n", err)
	}
}

