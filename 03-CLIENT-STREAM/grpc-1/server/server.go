package main

import (
	"context"
	proto "grpc-1/pb"
	datastore "grpc-1/util/dataStore"
	"grpc-1/util/fetcher"
	_"log"
)

type Server struct {
	proto.UnimplementedNewserviceServer
}

func (s *Server) GetNewsStream(req *proto.NewsRequest, stream proto.Newservice_GetNewsStreamServer) error {

	fetchedNews, err := fetcher.FetchNews("us")
	if err != nil {
		return err
	}

	go func() {

	}()

	for _, news := range fetchedNews {
		response := ArticleToNews(news)
		if err := stream.Send(response); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) GetNewsBulk(ctx context.Context, req *proto.NewsRequest) (*proto.BulkNews, error) {
	resp := &proto.BulkNews{}

	fetchedArticles, err := fetcher.FetchNews("en")
	// fetchedArticles, err := fetcher.FetchNewsTest("en")

	if err != nil {
		return nil, err
	}
	go datastore.SaveNewsDB(fetchedArticles)
	var newslist []*proto.News
	for _, article := range fetchedArticles {
		// log.Printf("data-fetched: ",article)
		// println("")
		newslist = append(newslist, ArticleToNews(article))
	}
	resp.News = newslist

	return resp, nil
}

func (s *Server) GetFreshNews(ctx context.Context, req *proto.NewsRequest) (*proto.BulkNews, error) {
	resp := &proto.BulkNews{}

	fetchedArticles, err := fetcher.FetchNews("en")

	if err != nil {
		return nil, err
	}
	go datastore.SaveNewsDB(fetchedArticles)
	var newslist []*proto.News
	for _, article := range fetchedArticles {
		// log.Printf("parsed article: ", ArticleToNews(article))
		// println("")
		newslist = append(newslist, ArticleToNews(article))
	}
	resp.News = newslist

	return resp, nil
}

func (s *Server) GetDBNews(ctx context.Context, req *proto.NewsRequest) (*proto.BulkNews, error) {
	resp := &proto.BulkNews{}
	dbArticle, err := datastore.GetDBNews()
	if err != nil {
		return nil, err
	}
	var newslist []*proto.News
	for _, article := range dbArticle {
		newslist = append(newslist, ArticleToNews(article))
	}
	resp.News = newslist
	return resp, nil
}
