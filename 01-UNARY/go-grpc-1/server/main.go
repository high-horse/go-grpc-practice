package main

import (
	"context"
	"fmt"
	proto "grpc1/helloGrpc"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{
	proto.UnimplementedHelloWorldServer
}

func main() {
	listener, tcpErr := net.Listen("tcp", ":9000")
	if tcpErr != nil {
		log.Fatalln(tcpErr)
	}
	serv := grpc.NewServer()
	proto.RegisterHelloWorldServer(serv, &server{})
	reflection.Register(serv)

	if err := serv.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}

func (s *server) SayHelloWorld(ctx context.Context, req *proto.HelloWorldRequest) (*proto.HelloWorldResponse, error) {
	fmt.Println("recieved req :", req.Name)
	return &proto.HelloWorldResponse{Message: "Hello " + req.Name}, nil
}
