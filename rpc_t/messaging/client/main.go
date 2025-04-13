package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "protos/message"
)

var (
	serverPort = flag.String("port", ":50051", "Port to listen on")
	fullName = flag.String("fullName", "John Doe", "Full name to send message.")
	// flag.int
	factorOf = flag.Int("factor", 120, "factor of number")
	f_name, l_name string
)

func parseName() {
	nameArr := strings.Split(*fullName, " ")
	f_name = nameArr[0]
	l_name = strings.Join(nameArr[1:], " ")
}


func main() {
	flag.Parse()
	parseName()

	conn, err := grpc.NewClient(*serverPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("could not create connection to the server ", err)
	}
	defer conn.Close()
	c := pb.NewMessageServiceClient(conn)

	doUnaryCall(c)
	go doServerStream(c)
	doServerStreamFactor(c)
}

func doUnaryCall(c pb.MessageServiceClient) {
	fmt.Println("initiating unary client")
	req := &pb.MessageUnaryReq{
		User: &pb.User{
			FirstName: f_name,
			LastName: l_name,
		},
	}
	res, err := c.MessageUnary(context.Background(), req)
	if err != nil {
		log.Fatal("error from unary client, ", err)
	}
	fmt.Println("response from unary message server :", res)
}

func doServerStream(c pb.MessageServiceClient){
	fmt.Println("initiating server stream call")

	req := &pb.MessaageServerStreamReq{
		User: &pb.User{
			FirstName: f_name,
			LastName: l_name,
		},
	}
	stream, err := c.MessageServerStream(context.Background(), req)
	if err != nil {
		log.Fatal("error calling from server stream ", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("error from the server ", err)
		}
		fmt.Println("response from the server ", res)
	}
}

func doServerStreamFactor(c pb.MessageServiceClient) {
	req := &pb.MessageServiceStreamFactorReq{
		Initial: int32(*factorOf),
	}
	factors := make([]int32, 0)
	stream , err := c.MessageServerStreamFactor(context.Background(), req)
	if err != nil {
		log.Fatal("error got from factor server stream " , err)
	}
	
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("error recieved from server stream factor ", err)
		}
		fmt.Println("recieved from server ", res)
		
		factors = append(factors, res.GetFactor())
	}
	strFactors := strings.Trim(strings.Replace(fmt.Sprint(factors), " ", ", ", -1), "[]")
	fmt.Printf("factors of %d are \n %s \n", *factorOf, strFactors)
}
