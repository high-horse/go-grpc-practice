package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"

	pb "protos/api/api"
)

var (
	port = flag.Int("port", 50051, "server port")
)

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal("could not start listener.")
	}
	defer listener.Close()
	
	s := grpc.NewServer()
	pb.RegisterGreetServiceServer(s, &GreetServer{})
	
	fmt.Println("starting server in port ", *port)
	if err := s.Serve(listener); err != nil {
		log.Fatal("could not start server")
	}
	
}

type GreetServer struct {
	pb.UnimplementedGreetServiceServer
}

func (s *GreetServer)  GreetMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageReply, error) {
	fmt.Println("recieved greet mssage request ", req)
	res := fmt.Sprintf("hello greetings %s %s", req.GetGreetingName().FirstName, req.GetGreetingName().LastName)
	
	return &pb.MessageReply{Result: res}, nil
}

func (s *GreetServer) GreetMessageManyServer(
	req *pb.GreetMessageManyServerReq, 
	stream grpc.ServerStreamingServer[pb.GreetMessageManyServerRes],
	// stream pb.GreetService_GreetMessageManyServerServer,
) error {
	fmt.Println("greet messaage many server stream rpc envoked with ", req)
	fName := req.GetGreetName().GetFirstName()
	lName := req.GetGreetName().GetLastName()
	// limit := 10
	for i := range 10 {
	// for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("Hello '%s %s'  -- counting %d", fName, lName, i)
		res := &pb.GreetMessageManyServerRes{
			Result: msg,
		}
		
		if err := stream.Send(res); err != nil {
			return err
		}
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*GreetServer) LongGreetClientStream(stream grpc.ClientStreamingServer[pb.LongGreetRequest, pb.LongGreetResponse]) error{
	fmt.Println("Long Greet Client Stream invoked.")
	res := make([]string, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("error from client stream, ", err)
		}
		res = append(res, fmt.Sprintf("hello %s %s", req.GetGreeting().GetFirstName(), req.GetGreeting().GetFirstName()))
	}
	return stream.SendAndClose(&pb.LongGreetResponse{Result: strings.Join(res, "\n")})
}

func (*GreetServer) GreetEveryone(stream grpc.BidiStreamingServer[pb.GreetEveryoneReq, pb.GreetEveryoneRes]) error {
	fmt.Println("GreetEveryone envoked")
	
	for {
		req, err := stream.Recv()
		if err == io.EOF {	
			break
		}
		if err != nil {
			return err
		}
		
		res := &pb.GreetEveryoneRes{Result: "helllo "+ req.GetGreeting().GetFirstName() + " !"}
		if err := stream.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (*GreetServer)	GreetWithDeadline(ctx context.Context, req *pb.MessageRequest) (*pb.MessageReply, error) {
	fmt.Println("GreetWithDeadline called with ", req.GetGreetingName())
	// return nil, status.Error(codes.Unimplemented, "method not implemented")
	
	res := fmt.Sprintf("hello greetings %s %s", req.GetGreetingName().FirstName, req.GetGreetingName().LastName)
	time.Sleep(3 * time.Second)
	if err := ctx.Err(); err != nil {
		log.Println("client cancelled the request.")
		return nil, status.Error(codes.Canceled, "client cancelled the request.")
	}
	return &pb.MessageReply{Result: res}, nil
}