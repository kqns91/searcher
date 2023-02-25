package repository

import (
	"context"

	"github.com/kqns91/searcher/api/model"
)

type OpenSearchRepository interface {
	Search(ctx context.Context, index []string, query string, from, size int) (*model.SearchResponse, error)
}
