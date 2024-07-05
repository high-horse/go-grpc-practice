package main

import (
	"context"
	proto "grpc1/helloGrpc"
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client proto.HelloWorldClient

func main() {
	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	client = proto.NewHelloWorldClient(conn)
	router := gin.Default()
	router.GET("/hello/:message", clientConnectionServer)
	router.Run(":8000")
}

func clientConnectionServer(c *gin.Context) {
	message := c.Param("message")

	req := &proto.HelloWorldRequest{Name: message}
	res, err := client.SayHelloWorld(context.Background(), req)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": res.GetMessage(),
	})
}