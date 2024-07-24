package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const URL = "https://newsapi.org/v2/top-headlines?sortBy=popularity&country=us"
const KEY = "&apiKey=bb67b0b40e014e9fb990a274637e73ad"

func FetchNews(country string) {
	url := URL + KEY
	res, err := http.Get(url)
	Check(err)
	defer res.Body.Close()
	allnews, err := io.ReadAll(res.Body)
	Check(err)

	var newsResponse HttpResponse
	err = json.Unmarshal(allnews, &newsResponse)
	for _, news := range newsResponse.Articles {
		// fmt.Printf("",news)
		fmt.Println(news.Source.Name)
		fmt.Println(news.Author)
		fmt.Println(news.Title)
	}

}

func Check(err error) {
	if err != nil {
		log.Fatalf("Error Encountered : ", err)
		os.Exit(1)
	}
}
