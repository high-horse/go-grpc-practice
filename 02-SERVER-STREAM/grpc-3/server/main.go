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

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	proto.UnimplementedProductServiceServer
}

const URL = "https://fakestoreapi.com/products"

func main() {
	// res, err := CallHttpRequest()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(res)

	// return

	address := ":50051"
	println("server starting at " + address)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Couldnot start server \n", err)
	}

	serv := grpc.NewServer()
	proto.RegisterProductServiceServer(serv, &Server{})

	reflection.Register(serv)

	if err := serv.Serve(listener); err != nil {
		log.Fatalf("Could not start server \n", err)
	}
}

func (s *Server) GetProduct(context.Context, *proto.ProductRequest) (*proto.ProductList, error) {
	var products proto.ProductList
	var err error

	productsSlice, err := CallHttpRequest()
	if err != nil {
		return nil, err
	}
	for _, product := range productsSlice.Products {
		products.Products = append(products.Products, product)
	}
	

	return &products, nil
}


func CallHttpRequest() (*proto.ProductList, error) {
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

    fmt.Printf("res: %s\n", string(body)) // Print the body for debugging

    var products []*proto.ProductResponse
    err = json.Unmarshal(body, &products)
    if err != nil {
        return nil, err
    }

    response := &proto.ProductList{Products: products}
    return response, nil
}
