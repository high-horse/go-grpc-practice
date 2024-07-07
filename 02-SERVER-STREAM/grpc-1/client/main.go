package main

import (
	"context"
	"grpc-1/pb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewWeatherServiceClient(conn)

	req := &pb.WeatherRequest{
		City: "London",
	}
	stream, err := c.GetWeatherUpdates(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to call GetWeatherUpdates: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Fatalf("error recieving: %v", err)
			break
		}
		log.Printf(
			"Weather Update: %s, %s, %s, %s", 
			res.City, res.Temperature, res.Weather, res.Update,
		)
	}
}