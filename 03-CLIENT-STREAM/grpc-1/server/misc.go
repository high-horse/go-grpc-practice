package main

import (
	"fmt"
	"grpc-1/util/fetcher"
	proto "grpc-1/pb"
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

func ArticleToNews(article fetcher.Article) *proto.News {
	return &proto.News{
		Source: &proto.Source{
			Id: article.Source.ID,
			Name: article.Source.Name,
		},
		Author: article.Author,
		Title: article.Title,
		Description: article.Description,
		Url: article.URL,
		PublishedAt: article.PublishedAt,
	}
}