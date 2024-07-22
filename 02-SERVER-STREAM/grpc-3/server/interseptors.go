package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

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

func streamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Printf("INtercepted stream call : %v", info.FullMethod)

	err := handler(srv, ss)
	if err != nil {
		log.Printf("Error from Stream call: %v", err)
		return status.Errorf(status.Code(err), "Stream call failed: %v", err)
	}
	return nil
}