package fetcher

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

const URL = "https://newsapi.org/v2/top-headlines?sortBy=popularity&country=us"
const KEY = "&apiKey=bb67b0b40e014e9fb990a274637e73ad"

func FetchNews(country string) ([]Article, error){
	url := URL + KEY
	res, err := http.Get(url)
	Check(err)
	defer res.Body.Close()
	allnews, err := io.ReadAll(res.Body)
	Check(err)

	var newsResponse *HttpResponse
	err = json.Unmarshal(allnews, &newsResponse)
	Check(err)

	return newsResponse.Articles, nil
}

func Check(err error) {
	if err != nil {
		log.Fatalf("Error Encountered : ", err)
		os.Exit(1)
	}
}
