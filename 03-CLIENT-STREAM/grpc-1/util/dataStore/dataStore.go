package datastore

import (
	"context"
	// "fmt"
	db "grpc-1/database"
	"grpc-1/util/fetcher"
	"log"
)



func SaveNewsDB(fetchedArticle []fetcher.Article) error {

	for _, article := range fetchedArticle {
		// continue
		
		source, news := ArticleToDBData(article)

		// fmt.Println(source)
		// fmt.Println(news)
		// continue
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
