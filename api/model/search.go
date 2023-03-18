package model

// request
type SearchTemplateRequest struct {
	ID     string           `json:"id"`
	Params BlogSearchParams `json:"params"`
}

type BlogSearchParams struct {
	Query []string `json:"query"`
	From  int      `json:"from"`
	Size  int      `json:"size"`
}

// response
type SearchResponse struct {
	Hits Hits `json:"hits,omitempty"`
}

type Hits struct {
	Total    Total   `json:"total,omitempty"`
	MaxScore float64 `json:"max_score,omitempty"`
	Hits     []Hit   `json:"hits,omitempty"`
}

type Total struct {
	Value    int    `json:"value,omitempty"`
	Relation string `json:"relation,omitempty"`
}

type Hit struct {
	Score     float64             `json:"_score,omitempty"`
	Source    Source              `json:"_source,omitempty"`
	Highlight map[string][]string `json:"highlight,omitempty"`
}

type Source struct {
	ArtiCode string   `json:"arti_code,omitempty"`
	Title    string   `json:"title,omitempty"`
	Member   string   `json:"member,omitempty"`
	Date     string   `json:"date,omitempty"`
	Link     string   `json:"link,omitempty"`
	Images   []string `json:"images,omitempty"`
	Comment1 string   `json:"comment1,omitempty"`
	KijiCode string   `json:"kijicode,omitempty"`
	Body     string   `json:"body,omitempty"`
}
