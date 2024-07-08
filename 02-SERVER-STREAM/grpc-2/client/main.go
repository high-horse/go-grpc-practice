package main

import (
	"context"
	"grpc-2/pb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error connecting : %v", err)
	}
	defer conn.Close()
	c := pb.NewProcessesClient(conn)

	req := &pb.ProcessRequest{}
	stream, err := c.GetProcessesInfo(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling GetProcessesInfo : %v", err)
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}
		log.Println(" Process Id:",resp.ProcessId)
		log.Println(" Process Name:",resp.ProcessName)
		log.Printf(" Cpu Usage: %.4f %",resp.CpuUsage)
		log.Printf(" Memory Usage: %.2f %s",resp.MemoryUsage, resp.MemUnit)
		println("")
	}
}