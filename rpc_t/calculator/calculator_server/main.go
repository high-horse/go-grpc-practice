package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"context"

	"google.golang.org/grpc"
	
	pb "protos/calculator/api"
)

var (
	port = flag.Int("port", 50051, "server port")
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal("failed to start listener.")
	}
	defer listener.Close()
	
	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &GServer{})
		log.Println("starting server")
	if err := s.Serve(listener); err != nil {
		log.Fatal("failed to serve")
	} 
}

type GServer struct {
	pb.UnimplementedCalculatorServiceServer
}

func (s *GServer) Sum(ctx context.Context, in *pb.CalculateSumRequest) (*pb.CalculateSumResponse, error) {
	fmt.Println("envoked sum wuth :", in)
	result := in.GetFirstNum() + in.GetSecondNum()
	fmt.Println("result found :", result)
	
	return &pb.CalculateSumResponse{ Result: result }, nil
}