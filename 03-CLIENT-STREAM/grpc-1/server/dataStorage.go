package main

// import (
// 	"database/sql"
// 	proto "grpc-1/pb"
// 	_ "grpc-1/store/database"
// )

// func StoreNews() {
	
// }


// // this implementation of generics mignt work, who knows??


// // Define a type constraint for the generic function
// type NewsRow interface {
// 	GetSource() string
// 	GetSourceName() string
// 	GetAuthor() sql.NullString
// 	GetTitle() sql.NullString
// 	GetDescription() sql.NullString
// 	GetPublishedAt() sql.NullTime
// }


// // Define the generic function
// func DbNewsToNews[T NewsRow](news T) *proto.News {
// 	var publishedAt string
// 	if news.GetPublishedAt().Valid {
// 		publishedAt = news.GetPublishedAt().Time.Format(time.RFC3339) // Convert time to string format
// 	}

// 	protoNews := &proto.News{
// 		Source: &proto.Source{
// 			Id:   news.GetSource(),
// 			Name: news.GetSourceName(),
// 		},
// 		Author:      news.GetAuthor().String,
// 		Title:       news.GetTitle().String,
// 		Description: news.GetDescription().String,
// 		PublishedAt: publishedAt,
// 	}
// 	return protoNews
// }



// /*

// // this implementation of generics is mistake
// func generic_circumference[r int | float32](radius r) {
 
//     c := 2 * 3 * radius
//     println("The circumference is: ", c) //comment
 
// }
// func DbNewsToNews [T storeDb.GetSourceBasedNewsRow | GetSingleNewsRow ](news T) *proto.News {
// 	protoNews := &proto.News{
// 		Source: &proto.Source{
// 			Id:   news.Source,
// 			Name: news.SourceName,
// 		},
// 		Author:      news.Author.String,
// 		Title:       news.Title.String,
// 		Description: news.Description.String,
// 		PublishedAt : news.Publishedat,
// 	}
// 	return protoNews
// }

//  */