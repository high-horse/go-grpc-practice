package main

import (
	proto "grpc-1/pb"
	storeDb "grpc-1/store/database"
)

func generic_circumference[r int | float32](radius r) {
 
    c := 2 * 3 * radius
    println("The circumference is: ", c) //comment
 
}
func DbNewsToNews [T storeDb.GetSourceBasedNewsRow | storeDb.GetSinglNewsRow](news T) *proto.News {
	protoNews := &proto.News{
		Source: &proto.Source{
			Id:   news.Source,
			Name: news.SourceName,
		},
		Author:      news.Author.String,
		Title:       news.Title.String,
		Description: news.Description.String,
		PublishedAt : news.Publishedat,
	}
	return protoNews
}