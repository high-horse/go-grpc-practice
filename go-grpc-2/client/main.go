package main

import (
	"context"
	"fmt"
	proto "grpc-2/invoice"
	"log"

	// "log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	log.Println("connecting at :50051...")
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := proto.NewInvoiceClient(conn)
	req := &proto.InvoiceRequest{
		Id: "test",
	}

	log.Println("requesting invoice...")
	res, err := client.GetInvoice(context.Background(), req)
	if err != nil {
		panic(err)
	}

	log.Println("got invoice...")
	fmt.Println(res)


	log.Println("creating invoice...")
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
	


	log.Println("updating invoice...")
	updInv := &proto.UpdateInvoiceRequest{
		Id: "test",
		PayerId: "test",
		PayerName: "test",
		PayerEmail: "test",
	}

	updRes, err := client.UpdateInvoice(context.Background(), updInv)
	if err != nil {
		panic(err)
	}
	fmt.Println("Updated Invoice :\n",updRes)


	log.Println("deleting invoice...")
	delRes, err := client.DeleteInvoice(context.Background(), &proto.DeleteInvoiceRequest{
		Id: "test",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Deleted Invoice :\n",delRes)

}
