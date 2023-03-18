package model

type Response struct {
	Total  int   `json:"total"`
	Result []any `json:"result"`
}

type Blog struct {
	ArtiCode  string   `json:"arti_code"`
	Title     string   `json:"title"`
	Member    string   `json:"member"`
	Date      string   `json:"date"`
	Link      string   `json:"link"`
	Images    []string `json:"images"`
	Highlight []string `json:"highlight"`
}

type Comment struct {
	Comment1  string   `json:"comment1"`
	Date      string   `json:"date"`
	KijiCode  string   `json:"kijicode"`
	Highlight []string `json:"highlight"`
}
