package datastore

import (
	"grpc-1/store/database"
	"grpc-1/util/fetcher"
	"database/sql"
	"time"
)

func ArticleToDBData(article fetcher.Article ) (database.CreateSourceParams, database.CreateNewsParams) {
	source := database.CreateSourceParams{
		SourceID: article.Source.ID,
		SourceName: article.Source.Name,
	}
	
	news := database.CreateNewsParams{
		SourceID: source.SourceID,
		Author: StringToNullString(article.Author),
		Title: StringToNullString(article.Title),
		Description: StringToNullString(article.Description),
		Publishedat: StringToNullTime(article.PublishedAt),
	}
	return source, news
}

func DBNewsToArticle(dbDews []database.GetAllNewsRow) []fetcher.Article {
	articles := []fetcher.Article{}
	
	for _, news := range dbDews {
		articles = append(articles, fetcher.Article{
			Source: fetcher.Source{
				ID: news.Source,
				Name: news.SourceName,
			},
			Author: NullStringToString(news.Author),
			Title: NullStringToString(news.Title),
			Description: NullStringToString(news.Description),
			URL: "",
			PublishedAt: NullTimeToString(news.Publishedat),
		})
	}
	
	return articles
}

func StringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NullStringToString(ns sql.NullString) string {
	if ns.Valid{
		return ns.String
	}
	return ""
}

func StringToNullTime(s string) sql.NullTime {
    if s == "" {
        return sql.NullTime{Valid: false}
    }
    t, err := time.Parse(time.RFC3339, s)
    if err != nil {
        return sql.NullTime{Valid: false}
    }
    return sql.NullTime{
        Time:  t,
        Valid: true,
    }
}

func NullTimeToString(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format(time.RFC3339)
	}
	return ""
}

