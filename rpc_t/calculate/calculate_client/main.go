package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"sync"

	// "os"
	"time"

	pb "protos/calculate"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 50051, "Server port")
)

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(fmt.Sprintf(":%d", *port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("could not create connection ", err)
	}
	defer conn.Close()

	c := pb.NewCalculateServiceClient(conn)
	// calculateSum(c)
	// callPrimeNumberWrap(c)
	// calculateAvg(c)
	handleMaxNumber(c)
}

func calculateSum(c pb.CalculateServiceClient) {
	fmt.Println("calling calculate sum.")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	req := &pb.CalculateSumRequest{FirstNum: 1092, SecondNum: 9987}
	res, err := c.CalculateSum(ctx, req)
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.DeadlineExceeded {
			log.Println("Request timed out ", err)
		} else {
			log.Println("error calculating sum ", err)
		}
		return
	}
	fmt.Printf("the sum of %v is %v \n", req, res.GetResult())
}
func callPrimeNumberWrap(c pb.CalculateServiceClient) {
	primeNumberDecomposition(c, -10)
}

func primeNumberDecomposition(c pb.CalculateServiceClient, num int32) {
	fmt.Println("calling prime number decomposition ")
	req := &pb.PrimeNumberDecompositionReq{Number: num}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	stream, err := c.PrimeNumberDecomposition(ctx, req)
	if fault := checkGrpcErr(err); fault {
		return
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		} else if fault := checkGrpcErr(err); fault {
			return
		}
		fmt.Println("response from server ", res.GetResult())
	}
}

func calculateAvg(c pb.CalculateServiceClient) {
	firstSet := []int32{1, 4, 5, 6, 7, 8}
	secondSet := []int32{10, 47, 53, 65, 77, 88}
	handleCalculateAvgRpcCall(c, firstSet)
	handleCalculateAvgRpcCall(c, secondSet)
}

func handleCalculateAvgRpcCall(c pb.CalculateServiceClient, numArr []int32) {
	fmt.Println("calling handleCalculateAvgRpcCall")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Println("error occured ", err)
		return
	}
	for _, val := range numArr {
		fmt.Println("sending ", val)
		req := &pb.ComputeAverageReq{Number: val}
		err := stream.Send(req)

		if err != nil {
			if fault := checkGrpcErr(err); fault {
				return
			}
		}
		time.Sleep(time.Second * 1)
	}
	res, err := stream.CloseAndRecv()
	if fault := checkGrpcErr(err); fault {
		return
	}
	fmt.Printf("the average of %v is %v \n", numArr, res.GetAverage())
}

func handleMaxNumber(c pb.CalculateServiceClient) {
	handleDuplexMaxRpcCall(c, []int32{1, 5, 6, 3, 2, 6})

	// handleDuplexMaxRpcCall(c, []int32{1, 5, 6, 3, 2, 6, 7, 9, 1, 2, 5, 8, 3, 4, 9, 8, 10})
}

func handleDuplexMaxRpcCall(c pb.CalculateServiceClient, nums []int32) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	stream, err := c.FindMaximum(ctx)
	if err != nil {
		return fmt.Errorf("failed go create stream : %w", err)
	}
	var (
		wg sync.WaitGroup
		errCh = make(chan error, 2)
		maxCh = make(chan int32, 1)
	)
	
	// start sender 
	wg.Add(1)
	go func(){
		defer wg.Done()
		for _, num := range nums {
			select{
				
			case <-ctx.Done():
				return
			default:
				if err := stream.Send(&pb.FindMaximumReq{Number: num}); err != nil {
					errCh <- fmt.Errorf("error sending number : %w", err)
					return
				}
				log.Println("sent: ", num)
			}
		}
		if err := stream.CloseSend(); err != nil {
			errCh <- fmt.Errorf("error closing send stream : %w", err)
		}
	}()
	
	// start recieving
	wg.Add(1)
	go func (){
		defer wg.Done()
		var currentMax int32
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				maxCh <- currentMax
				return
			}
			if err != nil {
				errCh <- fmt.Errorf("error recieving : %w", err)
				return
			}
			currentMax = res.GetResult()
			log.Println("recieved new max ", currentMax)
		}
	}()
	
	// wait for completinon or error
	go func(){
		wg.Wait()
		close(errCh)
		close(maxCh)
	}()
	
	select{
		case err := <-errCh:
			cancel()
			return err
			
		case max := <- maxCh:
			log.Printf("final max of %v is %d", nums, max)
			return nil
	}
}

func handleDuplexMaxRpcCall_old(c pb.CalculateServiceClient, nums []int32) {
	fmt.Println("handleDuplexMaxRpcCall envoked")
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		fmt.Println("error occured ", err)
		return
	}

	waitCh := make(chan struct{})
	// go routine to send
	go func() {
		fmt.Println("sending ", nums)
		for _, n := range nums {
			fmt.Println("sending ", n)
			if err := stream.Send(&pb.FindMaximumReq{Number: n}); err != nil {
				log.Println(err)
				break
			}
			time.Sleep(time.Second * 1)
		}
		stream.CloseSend()
	}()

	// go routine to recieve
	var max int32
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("error encountered ", err)
			}
			fmt.Println("recieved ", res.GetResult())
			max = res.GetResult()
		}
		close(waitCh)
	}()

	<-waitCh
	fmt.Printf("the maximum of %v is %d \n", nums, max)
}

func checkGrpcErr(err error) bool {
	if err != nil {
		st := status.Convert(err)

		switch st.Code() {
		case codes.DeadlineExceeded:
			log.Println("deadline Exceeded ", err)
			return true

		case codes.InvalidArgument:
			log.Println("Invalid Argument ", err)
			return true

		case codes.Aborted:
			log.Println("aborted ", err)
			return true

		case codes.NotFound:
			log.Println("endpoint not found ", err)
			return true

		default:
			log.Println("unknown error ", err)
			return true
		}
	}

	return false
}
