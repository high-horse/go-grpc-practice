package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "protos/calculator/api"
)

func main() {
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("could not connect to the grpc server.")
	}
	defer conn.Close()
	
	c := pb.NewCalculatorServiceClient(conn)
	
	req := &pb.CalculateSumRequest{FirstNum: 10, SecondNum: 20}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatal("error occured :", err)
	}
	fmt.Println("result found :", res.GetResult())
}