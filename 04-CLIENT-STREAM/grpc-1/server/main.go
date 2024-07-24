package main

import (
	"context"
	"fmt"
	"grpc-1/fetcher"
	"log"
	"net"

	proto "grpc-1/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const PORT = ":50051"
type Server struct {
	// proto.UnimplementedNewserviceServer
	proto.UnimplementedNewserviceServer
}

func main() {
	println("Server starting at : %s", PORT)
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal(err)
	}

	serv := grpc.NewServer()
	proto.RegisterNewserviceServer(serv, &Server{})
	reflection.Register(serv)

	if err := serv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}




func (s * Server) GetNewsStream(ctx context.Context, in *proto.NewsRequest) (proto.Newservice_GetNewsStreamClient, error) {
	return nil, nil
}

func fetch() {
	articles, err := fetcher.FetchNews("us")
	if err != nil {
		log.Fatalf("error :", err)
	}
	for _, article := range articles {
		fmt.Printf("title: %v \n", article.Title)
	}
}

