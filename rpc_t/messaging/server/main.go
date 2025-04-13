package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "protos/message"
)

var (
	port = flag.Int("port", 50051, "port to listen on")
)

func main() {
	flag.Parse()
	p := fmt.Sprintf(":%d", *port)
	fmt.Println("startng server on ", p)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	
	s := grpc.NewServer()
	pb.RegisterMessageServiceServer(s, &MessageServer{})
	s.Serve(listener)
}

type MessageServer struct {
	pb.UnimplementedMessageServiceServer
}

func (*MessageServer) MessageUnary(ctx context.Context,req *pb.MessageUnaryReq) (*pb.MessageUnaryRes, error){
	fmt.Println("MessageUnary initiated with ", req)
	greetMsg := fmt.Sprintf("hello %s %s", req.GetUser().GetFirstName(), req.GetUser().GetLastName())
	res := &pb.MessageUnaryRes{
		Result: greetMsg,
	}
	return res, nil
}

func (*MessageServer) MessageServerStream(req *pb.MessaageServerStreamReq, stream grpc.ServerStreamingServer[pb.MessageServerStreamRes]) error {
	preparedMsg := fmt.Sprintf("hello %s %s", req.GetUser().GetFirstName(), req.GetUser().GetLastName())
	
	for i := range 10 {
		res := &pb.MessageServerStreamRes{
			Result: fmt.Sprintf("%s %d", preparedMsg, i+1),
		}
		stream.Send(res)
		time.Sleep(1 * time.Second)
	}

	return nil
}

func (*MessageServer) MessageServerStreamFactor(req *pb.MessageServiceStreamFactorReq, stream grpc.ServerStreamingServer[pb.MessageServiceStreamFactorRes]) error {
	bigNum := req.GetInitial()
	var k int32 = 2
	for bigNum > 1 {
		if bigNum % k == 0 {
			if err := stream.Send(&pb.MessageServiceStreamFactorRes{
				Factor: k,
			}); err != nil {
				log.Fatal("error occured ", err)
			}
			bigNum = bigNum / k
		} else {
			k = k + 1
		}
		time.Sleep(time.Second * 1)
	}
	
	return nil
}