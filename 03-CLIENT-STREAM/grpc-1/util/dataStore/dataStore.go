package datastore

import (
	"context"
	"fmt"
	"grpc-1/store/database"
	"grpc-1/util/fetcher"
	"log"
)

var dbQueries database.Queries

func SaveNewsDB(fetchedArticle []fetcher.Article) error {

	for _, article := range fetchedArticle {
		// continue
		
		source, news := ArticleToDBData(article)

		fmt.Printf("%s \t %s\n", source.SourceID, news.Title)
		// continue
		_, err := dbQueries.CreateSource(context.Background(), source)
		if err != nil {
			log.Printf("Source DB err :%v", err)
		}

		_, err = dbQueries.CreateNews(context.Background(), news)
		if err != nil {
			log.Printf("News DB err :%v", err)
		}
	}
	return nil
}
