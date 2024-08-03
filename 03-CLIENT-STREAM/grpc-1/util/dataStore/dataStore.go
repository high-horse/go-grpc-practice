package datastore

import (
	"context"
	"grpc-1/store/database"
	"grpc-1/util/fetcher"
	"log"
)

var dbQueries *database.Queries

func SaveNewsDB(fetchedArticle *[]fetcher.Article) error{
	
	for _, article := range *fetchedArticle {
		source, news := ArticleToDBData(article)

		_, err := dbQueries.CreateSource(context.Background(), source)
		if err != nil {
			log.Printf("Source DB err :%v",err)
		}
		
		_, err = dbQueries.CreateNews(context.Background(),news)
		if err != nil {
			log.Printf("News DB err :%v",err)
		}
	}
	return nil
}
