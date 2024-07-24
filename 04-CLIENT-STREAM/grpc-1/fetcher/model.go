package fetcher

type HttpResponse struct {
	Status       string    `json:"status"`
	TotalResults uint32    `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Article struct {
	Source      Source `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string  `json:"url"`
	PublishedAt string `json:"publishedAt"`
	Content     string `json:"content"`
}

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BulkNews struct {
	news *[]Article
}
