package main

import (
	"grpc-1/pb"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterWeatherServiceServer(s, &server{})

	// Register reflectio service on grpc server
	reflection.Register(s)
	
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	pb.UnimplementedWeatherServiceServer
}

func (s * server) GetWeatherUpdates(req *pb.WeatherRequest, stream pb.WeatherService_GetWeatherUpdatesServer) error {
	for i := 0; i< 10 ; i++ {
		update := &pb.WeatherResponse {
			City : req.GetCity(),
			Weather: "sunny",
			Temperature: "30" + string(i),
			Update: "update #" + string(i),
		}
		if err := stream.Send(update); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}