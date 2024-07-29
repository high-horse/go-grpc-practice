package main

import (
	"context"
	"grpc-1/fetcher"
	proto "grpc-1/pb"
	"log"
)

type Server struct {
	proto.UnimplementedNewserviceServer
}


func (s *Server) GetNewsStream(req *proto.NewsRequest, stream proto.Newservice_GetNewsStreamServer) error {

	fetchedNews, err := fetcher.FetchNews("us")
	if err != nil {
		return err
	}

	for _, news := range fetchedNews{
		response := ArticleToNews(news)
		// source := proto.Source{
		// 	Id: news.Source.ID,
		// 	Name: news.Source.Name,
		// }
		// resNews := &proto.News{
		// 	Source: &source,

		// }
		if err := stream.Send(response); err != nil {
			return err
		}
	}

	return  nil
}


func (s *Server) GetNewsBulk(ctx context.Context, req *proto.NewsRequest) (*proto.BulkNews, error){
	resp := &proto.BulkNews{}

	fetchedArticles, err := fetcher.FetchNews("en")
	if err != nil {
		return nil, err
	}
	var newslist  []*proto.News
	for _, article := range fetchedArticles{
		log.Printf("data-fetched: ",article)
		println("")
		newslist = append(newslist, ArticleToNews(article))
	}
	resp.News = newslist

	return resp, nil
}
