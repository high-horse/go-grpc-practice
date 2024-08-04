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

func FetchNewsTest(country string) ([]Article, error) {
	var articles []Article
	a := Article{
		Source: Source{
                ID: "the-washington-post",
                Name: "The Washington Post",
            },
            Author: "Jo-Ann Finkelstein",
            Title: "Perspective | Gross and embarrassing — teen girls’ misconceptions about their periods - The Washington Post",
            Description: "Parents can help dispel confusion and misinformation and alleviate the anxiety and shame some girls may feel about menstruation.",
            PublishedAt: "2024-08-03T12:42:16Z",
	}
	articles = append(articles, a)
	
	return articles, nil
}

func Check(err error) {
	if err != nil {
		log.Fatalf("Error Encountered : ", err)
		os.Exit(1)
	}
}
