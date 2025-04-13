package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
	// "time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "protos/api/api"
)


var (
	addr = flag.String("addr", "localhost:50051", "adress to the server")
	name = flag.String("name", "hello some default name ", "greeter name")
)

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error creating new grpc client")
	}
	defer conn.Close()
	
	clent := pb.NewGreetServiceClient(conn)

	// doUnary(clent)
	doServerStream(clent)
}

func doUnary(clent pb.GreetServiceClient){
	fmt.Println("starting unary services")
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	
	names := strings.Split(*name, " ")
	greetname := &pb.GreetingName{
		FirstName: names[0],
		LastName: strings.Join(names[1:], " "),
	}
	r, err := clent.GreetMessage(context.Background(), &pb.MessageRequest{GreetingName: greetname})
	if err != nil {
		log.Fatal("could not recieve expected response from the server.")
	}
	
	fmt.Println(r)
}

func doServerStream(client pb.GreetServiceClient) {
	fmt.Println("server stream client envoked")
	
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	names := strings.Split(*name, " ")
	in := &pb.GreetMessageManyServerReq{
		GreetName: &pb.GreetingName{
			FirstName: names[0],
			LastName: strings.Join(names[1:], " "),
		},
	}
	recvStream, err := client.GreetMessageManyServer(context.Background(), in)
	if err != nil {
		log.Fatal("error occured while creating recvStream ", err)
	}
	
	for {
		
		msg, err := recvStream.Recv()
		if err == io.EOF  {
			break
		}
		if err != nil {
			log.Fatal("error while reading from the stream ", err)
		}
		
		fmt.Println("response from stream :", msg.GetResult())
	}
}