package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	opensearchv2 "github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"github.com/kqns91/searcher/api/model"
	"github.com/kqns91/searcher/api/repository"
)

type osearch struct {
	client *opensearchv2.Client
}

func New(client *opensearchv2.Client) repository.OpenSearchRepository {
	return &osearch{
		client: client,
	}
}

func (o *osearch) Search(ctx context.Context, index []string, query string, from, size int) (*model.SearchResponse, error) {
	words := regexp.MustCompile(` |　`).Split(query, -1)

	body, err := json.Marshal(model.SearchTemplateRequest{
		ID: "blog_search",
		Params: model.BlogSearchParams{
			Query: words,
			From:  from,
			Size:  size,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}

	req := opensearchapi.SearchTemplateRequest{
		Index: index,
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, o.client)
	if err != nil {
		return nil, fmt.Errorf("failed to search document: %w", err)
	}

	defer res.Body.Close()

	var v model.SearchResponse

	err = json.NewDecoder(res.Body).Decode(&v)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}

	return &v, nil
}
