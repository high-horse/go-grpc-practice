package main

import (
	"context"
	"fmt"
	"log"

	proto "grpc-1/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const SERVER = "localhost:50051"

func main() {
	println("server starting at : ", SERVER)
	conn, err := grpc.NewClient(SERVER, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("ERROR :", err)
	}
	defer conn.Close()

	client := proto.NewNewserviceClient(conn)

	getNewsBulk(client)
	// getNewsStream(client)
}

func getNewsBulk(client proto.NewserviceClient) {
	req := &proto.NewsRequest{}

	res, err := client.GetNewsBulk(context.Background(), req)
	if err != nil {
		log.Println("Error :", err)
	}
	fmt.Println("Response from the server: ")
	for _, result := range res.News {
		fmt.Printf("%s \t %s \n", result.Author, result.Description)
		println()
	}
}
