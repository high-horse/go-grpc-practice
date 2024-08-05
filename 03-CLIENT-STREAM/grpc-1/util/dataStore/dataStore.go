package datastore

import (
	"context"
	"fmt"
	// "fmt"
	db "grpc-1/database"
	"grpc-1/util/fetcher"
	"log"
)



func SaveNewsDB(fetchedArticle []fetcher.Article) error {

	for _, article := range fetchedArticle {
		source, news := ArticleToDBData(article)
		_, err := db.Queries.CreateSource(context.Background(), source)
		if err != nil {
			log.Printf("Source DB err :%v", err)
		}

		_, err = db.Queries.CreateNews(context.Background(), news)
		if err != nil {
			log.Printf("News DB err :%v", err)
		}
	}
	return nil
}

func GetDBNews() ([]fetcher.Article, error) {
	newsDB, err  := db.Queries.GetAllNews(context.Background())
	if err != nil{
		return nil, fmt.Errorf("GetDBNews : %v",err)
	}
	return DBNewsToArticle(newsDB), nil
}
