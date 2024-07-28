package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryInterseptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error){
	log.Printf("intersepted unary call : %v \n", info.FullMethod)

	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("error from unary interceptor : %v\n", err)
		return nil, status.Errorf(status.Code(err), "unary call failed, %v", err)
	}

	return resp, err
}

func StreamInterseptor (
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Printf("INtersepted stream call : %v", info.FullMethod)

	err := handler(srv, ss)
	if err != nil {
		log.Printf("Error from stream call : %v", err)
		return status.Errorf(status.Code(err), "strem call failed: %v", err)
	}
	return nil
}