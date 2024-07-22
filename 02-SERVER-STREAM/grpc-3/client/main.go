package main

import (
	"context"
	"fmt"
	"log"
	"os"

	proto "grpc-3/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/credentials/insecure"
)

const SERVER = "localhost:50051"

func unaryInterseptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("intercepted Unary call: %v", info.FullMethod)

	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("error from unary call : %v", err)
		return nil, status.Errorf(status.Code(err), "Unary call failed, %v", err)
	}

	

	return resp, err
}


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

	fmt.Printf("request from the server:")
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