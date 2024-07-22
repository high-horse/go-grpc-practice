package main

import (
	"context"
	"fmt"
	"log"
	"os"

	proto "grpc-3/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const SERVER = "localhost:50051"


func main() {
	log.Printf("Starting Server at %s \n", SERVER)
	conn, err := grpc.NewClient(SERVER, grpc.WithTransportCredentials(insecure.NewCredentials()))
	check(err)

	client := proto.NewProductServiceClient(conn)
	getSingleProduct(client)
}

func getSingleProduct(client proto.ProductServiceClient) {
	req := &proto.ProductRequest{}

	res, err := client.GetProduct(context.Background(), req)
	check(err)

	fmt.Printf("response from the server:")
	for _, result := range res.Products {
		fmt.Printf("id :%d \t, title :%s \n", result.Id, result.Title)
	}
}

func check(err error) {
	if err != nil {
		println("error :", err)
		os.Exit(1)
	}
}