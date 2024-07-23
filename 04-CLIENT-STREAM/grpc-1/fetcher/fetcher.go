package fetcher

import (
	"log"
	"net/http"
	"os"
)

const URL = "https://newsapi.org/v2/top-headlines?sortBy=popularity&country=us"
const KEY = "&apiKey=bb67b0b40e014e9fb990a274637e73ad"

func FetchNews(country string, ) {
	url := URL + KEY
	_, err := http.Get(url)
	Check(err)
}

func Check(err error) {
	if err != nil {
		log.Fatalf("Error Encountered : ", err)
		os.Exit(1)
	}
}