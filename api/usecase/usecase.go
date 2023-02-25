package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kqns91/searcher/api/model"
	"github.com/kqns91/searcher/api/repository"
)

type Usecase interface {
	Search(ctx context.Context, query string, from, size string) (*model.SearchResponse, error)
}

type ucase struct {
	search repository.OpenSearchRepository
}

func New(repo repository.OpenSearchRepository) Usecase {
	return &ucase{
		search: repo,
	}
}

func (u *ucase) Search(ctx context.Context, query string, from, size string) (*model.SearchResponse, error) {
	if query == "" {
		return &model.SearchResponse{}, nil
	}

	var err error
	f := 0
	s := 30

	if from != "" {
		f, err = strconv.Atoi(from)
		if err != nil {
			return nil, fmt.Errorf("failed to convert from: %w", err)
		}
	}

	if size != "" {
		s, err = strconv.Atoi(size)
		if err != nil {
			return nil, fmt.Errorf("failed to convert size: %w", err)
		}
	}

	res, err := u.search.Search(ctx, []string{"blogs"}, query, f, s)
	if err != nil {
		return nil, fmt.Errorf("failed to search document: %w", err)
	}

	return res, nil
}
