package main

import (
	"context"
	"encoding/json"
	"fmt"
	proto "grpc-3/pb"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	proto.UnimplementedProductServiceServer
}

const URL = "https://fakestoreapi.com/products"
const ADDRESS = ":50051"


func unaryInterseptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("intercepted Unary call: %v", info.FullMethod)

	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("error from unary call : %v", err)
		return nil, status.Errorf(status.Code(err), "Unary call failed, %v", err)
	}

	

	return resp, err
}


func main() {
	// res, err := CallHttpRequest()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, result := range res {
	// 	fmt.Printf("id : %d \t title :%s \n", result.Id, result.Title)
	// }

	// return

	
	println("server starting at " + ADDRESS)

	listener, err := net.Listen("tcp", ADDRESS)
	if err != nil {
		log.Fatalf("Couldnot start server \n", err)
	}

	serv := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterseptor),
	)
	proto.RegisterProductServiceServer(serv, &Server{})

	reflection.Register(serv)

	if err := serv.Serve(listener); err != nil {
		log.Fatalf("Could not start server \n", err)
	}
}

func (s *Server) GetProduct(context.Context, *proto.ProductRequest) (*proto.ProductList, error) {
	// return &proto.ProductList{}, nil
	var products proto.ProductList
	var err error

	// productsSlice, err := CallHttpRequest()
	productSl, err := CallHttpRequest()
	if err != nil {
		return nil, err
	}

	for _, prod := range productSl{
		products.Products = append(products.Products, prod)
	}
	
	if err != nil {
		return nil, err
	}

	return &products, nil
}

func (s *Server) GetProductStream(req *proto.ProductRequest, stream proto.ProductService_GetProductStreamServer) error {
	for {
		httpRes, err := CallHttpRequest()
		if err != nil {
			return err
			break
		}
		for _, res := range httpRes {
			if err := stream.Send(res); err != nil {
				return err
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
	return nil
}

func CallHttpRequest() ([]*proto.ProductResponse, error) {
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("res: %s\n", string(body)) // Print the body for debugging

	var products []*proto.ProductResponse
	err = json.Unmarshal(body, &products)
	if err != nil {
		return nil, err
	}

	// response := &proto.ProductList{Products: products}
	return products, nil
}
