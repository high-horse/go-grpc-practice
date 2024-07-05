package main

import (
	"context"
	proto "grpc-2/invoice"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	proto.UnimplementedInvoiceServer
}

func main() {
	println("Server starting at :50051...")
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	serv := grpc.NewServer()
	proto.RegisterInvoiceServer(serv, &Server{})
	reflection.Register(serv)

	if err := serv.Serve(listener); err != nil {
		panic(err)
	}
}

func (s *Server) GetInvoice(ctx context.Context, req *proto.InvoiceRequest) (*proto.InvoiceResponse, error) {
	log.Println("new request")
	res := &proto.InvoiceResponse{
		Id:         req.Id,
		PayerId:    "test",
		PayerName:  "req.PayerName",
		PayerEmail: "req.PayerEmail",
	}
	return res, nil
}

func (s *Server) CreateInvoice(ctx context.Context, req *proto.CreateInvoiceRequest) (*proto.CreateInvoiceResponse, error) {
	log.Println("new request")
	res := &proto.CreateInvoiceResponse{
		InvMessage: "test",
		Id:         req.Id,
		PayerId:    req.PayerId,
		PayerName:  req.PayerName,
		PayerEmail: req.PayerEmail,
	}
	return res, nil
}
