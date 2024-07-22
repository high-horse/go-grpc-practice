package main

import (
	"context"
	proto "grpc-3/pb"
	"time"
)

type Server struct {
	proto.UnimplementedProductServiceServer
}

func (s *Server) GetProduct(context.Context, *proto.ProductRequest) (*proto.ProductList, error) {
	var products proto.ProductList
	var err error

	productSl, err := CallHttpRequest()
	if err != nil {
		return nil, err
	}

	// for _, prod := range productSl{
	// 	products.Products = append(products.Products, prod)
	// }
	products.Products = append(products.Products, productSl...)
	
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
