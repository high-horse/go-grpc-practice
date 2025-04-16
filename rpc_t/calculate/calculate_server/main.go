package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "protos/calculate"

	"google.golang.org/grpc"
)


var (
	port = flag.Int("port", 50051, "Port for server")
)

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	defer listener.Close()
	if err != nil {
		log.Fatal("error creating listener ",err)
	}
	
	s := grpc.NewServer()
	fmt.Println("starting server in ", *port)
	pb.RegisterCalculateServiceServer(s, &CalculateServer{})
	if err := s.Serve(listener); err != nil {
		log.Fatal("error binding port to server ", err)
	}
	
}

type CalculateServer struct {
	pb.UnimplementedCalculateServiceServer
}

func (*CalculateServer) CalculateSum(ctx context.Context, req *pb.CalculateSumRequest) (*pb.CalculateSumResponse, error) {
	fmt.Println("Calculate sun invoked with ", req)
	firstNum := req.GetFirstNum()
	secondNum := req.GetSecondNum()
	
	sum := firstNum + secondNum
	time.Sleep(5 * time.Second)
	return &pb.CalculateSumResponse{Result: sum}, nil
}