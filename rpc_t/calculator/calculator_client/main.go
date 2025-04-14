package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "protos/calculator/api"
)

func main() {
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("could not connect to the grpc server.")
	}
	defer conn.Close()
	
	c := pb.NewCalculatorServiceClient(conn)
	
	// doUnaryCall(c)
	// doServerStream(c)
	// doClientStream(c)
	// doBiDiStream(c)
	doUnaryErrorHandling(c)
}

func doUnaryCall(c pb.CalculatorServiceClient) {
	req := &pb.CalculateSumRequest{FirstNum: 10, SecondNum: 20}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatal("error occured :", err)
	}
	fmt.Println("result found :", res.GetResult())
}

func doServerStream(c pb.CalculatorServiceClient) {
	factors := make([]int64, 0)
	number := 120
	req := &pb.PrimeNumberDecompositionReq{Number: int64(number)}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatal("error while creating server stream ", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error occured during stream ", err)
			return
		}
		log.Println("recieved response ", res)
		factors = append(factors, res.GetFactor())
	}
	
	fmt.Println("the factors of " +  strconv.Itoa(number) + " is ",factors)
}

func doClientStream(c pb.CalculatorServiceClient) {
	nums := []int64 {1,2,3,4}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatal("error connecting to the client stream server")
	}
	
	for _, n := range nums {
		fmt.Println("Sending ", n)
		if err := stream.Send(&pb.ComputeAverageReq{Number: n}); err != nil {
			log.Println("Error while sending throuh stream ", err)
		}
		time.Sleep(time.Second * 1)
	}
	
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("error encountered while closing stream, " , err)
	}
	fmt.Printf("Average of %v is %.2f\n", nums, res.GetAverage())
}

func doBiDiStream(c pb.CalculatorServiceClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatal("error occired while creating stream ",err)
	}
	
	nums := []int32{1,5,3,6,2,20}
	waitCh := make(chan struct{})
	var max int32
	go func(){
		for _, n := range nums {
			fmt.Println("sending ", n)
			if err := stream.Send(&pb.FindMaximumReq{Number: n}); err != nil {
				log.Fatal("error occured while senging stream ", err)
			}
			time.Sleep(1*time.Second)
		}
		_ = stream.CloseSend()
	}()
	
	go func() {
		for {
			res, err := stream.Recv()
			if err ==io.EOF {
				break
			}
			if err != nil {
				log.Fatal("error occured while recieving stream, ", err)
			}
			max = res.GetMaximum()
			fmt.Println("recieved ", res.GetMaximum())
		}
		close(waitCh)
	}()
	
	<-waitCh
	fmt.Printf("maximum of %v is %d\n", nums, max)
	// fmt.Fprintf("maximim of %v is %d", (nums), max)
}

func doUnaryErrorHandling(c pb.CalculatorServiceClient) {
	callSqRoot(c, 100)
	callSqRoot(c, 10)
	callSqRoot(c, -1)
	callSqRoot(c, 0)
	
}

func callSqRoot(c pb.CalculatorServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &pb.SquareRootReq{Number: number})
	if err != nil {
		errCode := status.Code(err)
		switch errCode {
			case codes.InvalidArgument:
				fmt.Printf("Error: invalid argument -%v\n", err)
				
			case codes.Internal:
				fmt.Printf("Internal server error - %v\n", err)
				
			case codes.DeadlineExceeded:
				fmt.Printf("Error: Request timeout - %v\n", err)
				
			case codes.NotFound:
				fmt.Printf("Error: Not Found - %v\n", err)
				
			default:
				fmt.Printf("Unknown error: %v\n", err)
		}
		return
	}
	fmt.Printf("The square root of %v is %v \n", number, res.GetSqrootNumber())
}