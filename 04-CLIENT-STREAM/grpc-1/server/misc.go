package main

import (
	"fmt"
	"grpc-1/fetcher"
	"log"
)

func fetch() {
	articles, err := fetcher.FetchNews("us")
	if err != nil {
		log.Fatalf("error :", err)
	}
	for _, article := range articles {
		fmt.Printf("title: %v \n", article.Title)
	}
}

