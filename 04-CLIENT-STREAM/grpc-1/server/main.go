package main

import (
	"log"
	"net"

	proto "grpc-1/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const PORT = ":50051"
type Server struct {
	proto.UnimplementedNewserviceServer
}

func main() {
	println("Server starting at ", PORT)
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal(err)
	}

	// create a new grpc server
	serv := grpc.NewServer()
	proto.RegisterNewserviceServer(serv, &Server{})
	reflection.Register(serv)

	if err := serv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}


func (s *Server) GetNewsStream(req *proto.NewsRequest, stream proto.Newservice_GetNewsStreamServer) error {
	
	return  nil
}


