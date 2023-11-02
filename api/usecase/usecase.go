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
	Search(ctx context.Context, index string, query string, from, size string) (*model.Response, error)
}

type ucase struct {
	search repository.OpenSearchRepository
}

func New(repo repository.OpenSearchRepository) Usecase {
	return &ucase{
		search: repo,
	}
}

func (u *ucase) Search(ctx context.Context, index string, query string, from, size string) (*model.Response, error) {
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

	switch index {
	case "blogs", "comments":
	default:
		return &model.Response{}, nil
	}

	// commentsも検索できるようにする。
	sr, err := u.search.Search(ctx, []string{index}, query, f, s)
	if err != nil {
		return nil, fmt.Errorf("failed to search document: %w", err)
	}

	result := []any{}

	for _, h := range sr.Hits.Hits {
		highlight := []string{}

		keys := make([]string, 0, len(h.Highlight))

		for k := range h.Highlight {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			highlight = append(highlight, h.Highlight[k]...)
		}

		switch index {
		case "blogs":
			result = append(result, &model.Blog{
				ID:        h.ID,
				ArtiCode:  h.Source.ArtiCode,
				Title:     h.Source.Title,
				Member:    h.Source.Member,
				Date:      h.Source.Date,
				Link:      h.Source.Link,
				Images:    h.Source.Images,
				Highlight: highlight,
			})
		case "comments":
			result = append(result, &model.Comment{
				ID:        h.ID,
				Comment1:  h.Source.Comment1,
				Date:      h.Source.Date,
				KijiCode:  h.Source.KijiCode,
				Body:      h.Source.Body,
				Highlight: highlight,
			})
		}
	}

	res := &model.Response{
		Total:  sr.Hits.Total.Value,
		Result: result,
	}

	return res, nil
}
