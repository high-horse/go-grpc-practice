package fetcher

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

// const URL = "https://newsapi.org/v2/top-headlines?sortBy=popularity&country=np"
// const URL = "https://newsapi.org/v2/top-headlines?sortBy=newest&country=np"
// const KEY = "&apiKey=bb67b0b40e014e9fb990a274637e73ad"
const URL = "https://newsapi.org/v2/top-headlines?sortBy=newest&apiKey=bb67b0b40e014e9fb990a274637e73ad&country=us"

func FetchNews(country string) ([]Article, error) {
	// url := URL + KEY
	// res, err := http.Get(url)
	res, err := http.Get(URL)
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
