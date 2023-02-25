package model

type Response struct {
	Total int     `json:"total"`
	Blogs []*Blog `json:"blogs"`
}

type Blog struct {
	Title     string   `json:"title"`
	Member    string   `json:"member"`
	Created   string   `json:"created"`
	URL       string   `json:"url"`
	Highlight []string `json:"highlight"`
}
