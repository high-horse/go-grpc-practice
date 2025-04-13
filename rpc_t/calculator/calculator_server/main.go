package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"

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

func (*GServer) PrimeNumberDecomposition(req *pb.PrimeNumberDecompositionReq, stream grpc.ServerStreamingServer[pb.PrimeNumberDecompositionRes]) error {
	log.Println("PrimeNumberDecomposition envoked with  ", req)
	var k int64 = 2
	n := req.GetNumber()
	for n > 1 {
		if n % k == 0 {
			log.Println("sending response...", k)
			if err := stream.Send(&pb.PrimeNumberDecompositionRes{Factor: k}); err != nil {
				log.Println("error ", err)
			}
			n = n/k
		} else {
			k = k+1
		}
		time.Sleep(time.Second * 1)
	}	
	return nil
}

func (*GServer) ComputeAverage(stream grpc.ClientStreamingServer[pb.ComputeAverageReq, pb.ComputeAverageRes]) error {
	var sum float64 = 0
	var count float64 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.ComputeAverageRes {Average: (sum/count)})
		}
		if err != nil {
			log.Println("error occured during streaming ")
			continue
		}
		sum += float64(req.GetNumber())
		count ++
	}
}
