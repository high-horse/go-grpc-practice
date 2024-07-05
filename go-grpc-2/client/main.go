package main

import (
	"context"
	"fmt"
	proto "grpc-2/invoice"
	// "log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := proto.NewInvoiceClient(conn)
	req := &proto.InvoiceRequest{
		Id: "test",
	}

	res, err := client.GetInvoice(context.Background(), req)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)

	inv, err := client.CreateInvoice(context.Background(), &proto.CreateInvoiceRequest{
		Id: "test",
		PayerId: "test",
		PayerName: "test",
		PayerEmail: "test",
		Amount: "test",
		Currency: "test",
		Description: "test",
		Metadata: "test",
	})

	fmt.Println("Generated Invoice :\n",inv)
	
}