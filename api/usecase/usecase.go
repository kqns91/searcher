package usecase

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/kqns91/searcher/api/model"
	"github.com/kqns91/searcher/api/repository"
)

type Usecase interface {
	Search(ctx context.Context, query string, from, size string) (*model.Response, error)
}

type ucase struct {
	search repository.OpenSearchRepository
}

func New(repo repository.OpenSearchRepository) Usecase {
	return &ucase{
		search: repo,
	}
}

func (u *ucase) Search(ctx context.Context, query string, from, size string) (*model.Response, error) {
	if query == "" {
		return &model.Response{}, nil
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

	sr, err := u.search.Search(ctx, []string{"blogs"}, query, f, s)
	if err != nil {
		return nil, fmt.Errorf("failed to search document: %w", err)
	}

	blogs := []*model.Blog{}

	for _, h := range sr.Hits.Hits {
		highlight := []string{}

		keys := make([]string, 0, len(h.Highlight))
		for k := range h.Highlight {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		for _, k := range keys {
			highlight = append(highlight, h.Highlight[k]...)
		}

		blogs = append(blogs, &model.Blog{
			Title:     h.Source.Title,
			Member:    h.Source.Member,
			Created:   h.Source.Created,
			URL:       h.Source.URL,
			Highlight: highlight,
		})
	}

	res := &model.Response{
		Total: sr.Hits.Total.Value,
		Blogs: blogs,
	}

	return res, nil
}
