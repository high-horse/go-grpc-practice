package main

import (
	"context"
	proto "grpc1/helloGrpc"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	client := proto.NewHelloWorldClient(conn)
	req := &proto.HelloWorldRequest{Name: "Mark"}
	res, err := client.SayHelloWorld(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.GetMessage())
}