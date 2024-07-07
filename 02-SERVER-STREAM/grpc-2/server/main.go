package main

import (
	"grpc-2/pb"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")	
	if err != nil {
		log.Fatalf("Error listening : %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterProcessesServer(s, &server{})

	reflection.Register(s)

	if err  := s.Serve(listener); err != nil {
		log.Fatalf("Error serving : %v", err)
	}
}

