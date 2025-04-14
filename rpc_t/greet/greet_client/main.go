package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

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
	// doServerStream(clent)
	// doClientStreaming(clent)
	doBiDirectionStreaming(clent)
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

func doClientStreaming(c pb.GreetServiceClient) {
	requests := []*pb.LongGreetRequest{
		&pb.LongGreetRequest{Greeting : &pb.GreetingName{FirstName: "jonh",}},
		&pb.LongGreetRequest{Greeting : &pb.GreetingName{FirstName: "mike",}},
		&pb.LongGreetRequest{Greeting : &pb.GreetingName{FirstName: "jonny",}},
		&pb.LongGreetRequest{Greeting : &pb.GreetingName{FirstName: "sam",}},
		&pb.LongGreetRequest{Greeting : &pb.GreetingName{FirstName: "harry",}},
		
	}
	stream, err := c.LongGreetClientStream(context.Background())
	if err != nil {
		log.Fatal("error while creating stream for client stream.")
	}
	
	for _, req := range requests {
		log.Printf("sending request %v\n", req)
		if err := stream.Send(req); err != nil {
			log.Println("error encountered while sending stream, " , err)
		}
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("error encountered while closing stream, " , err)
		return
	}
	fmt.Printf("recieved from server \n%s\n" , res.GetResult())
	// LongGreetClientStream(context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[LongGreetRequest, LongGreetResponse], error)
}

func doBiDirectionStreaming(c pb.GreetServiceClient) {
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatal("error occured while creating stream, ", err)
	}
	
	waitCh := make(chan struct{})
	requests := []*pb.GreetEveryoneReq{
		&pb.GreetEveryoneReq{Greeting : &pb.GreetingName{FirstName: "jonh",}},
		&pb.GreetEveryoneReq{Greeting : &pb.GreetingName{FirstName: "mike",}},
		&pb.GreetEveryoneReq{Greeting : &pb.GreetingName{FirstName: "jonny",}},
		&pb.GreetEveryoneReq{Greeting : &pb.GreetingName{FirstName: "sam",}},
		&pb.GreetEveryoneReq{Greeting : &pb.GreetingName{FirstName: "harry",}},
		
	}
	go func(){
		for _, req := range requests {
			log.Println("sending ", req)
			if err := stream.Send(req); err != nil {
				log.Println("error recieved when sending stream ", err)
			}
			time.Sleep(1 * time.Second)
		}
		_ = stream.CloseSend()
	}()
	
	go func(){
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("error occured while recieving from stream ", err)
				break
			}
			
			fmt.Println("revieved from stream ", res.GetResult())
		}
		close(waitCh)
		
	}()
	
	<-waitCh
}