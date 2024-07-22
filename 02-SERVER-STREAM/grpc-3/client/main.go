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
	defer conn.Close()

	client := proto.NewProductServiceClient(conn)
	getSingleProduct(client)
	stream_products(client)
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

func stream_products(client proto.ProductServiceClient) {
	println("starting server streaming...")
	req := &proto.ProductRequest{}
	stream, err := client.GetProductStream(context.Background(), req)
	check(err)

	for {
		res, err := stream.Recv()
		check(err)

		fmt.Printf("product id : %d,\t title: %s \n", res.Id, res.Title)
	}
}

func check(err error) {
	if err != nil {
		println("error :", err)
		os.Exit(1)
	}
}