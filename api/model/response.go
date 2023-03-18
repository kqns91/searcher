package model

type Response struct {
	Total int     `json:"total"`
	Blogs []*Blog `json:"blogs"`
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
