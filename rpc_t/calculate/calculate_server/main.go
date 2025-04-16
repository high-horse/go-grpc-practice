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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*CalculateServer) PrimeNumberDecomposition(req *pb.PrimeNumberDecompositionReq, stream grpc.ServerStreamingServer[pb.PrimeNumberDecompositionRes]) error {
	fmt.Println("invoked Prime Number Decomposition with ", req)
	num := req.GetNumber()
	if num < 0 {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("number must be grater than 0; got: %d", num))
	}
	k := int32(2)
	for num > 1 {
		if err := stream.Context().Err(); err != nil {
			fmt.Println("client disconnected, stopping function ")
			return nil
		}
		if num % k == 0 {
			stream.Send(&pb.PrimeNumberDecompositionRes{Result: k})
			num = num / k
		} else {
			k ++
		}
		time.Sleep(time.Second * 2)
		fmt.Println("number reached ", k)
	}
	return nil
}