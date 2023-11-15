package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kqns91/searcher/api/usecase"
)

type Handler interface {
	SearchBlogs(ctx context.Context, c *gin.Context) error
	ListBlogs(ctx context.Context, c *gin.Context) error
	SearchComments(ctx context.Context, c *gin.Context) error
}

type httpHandler struct {
	uc usecase.Usecase
}

func New(uc usecase.Usecase) Handler {
	return &httpHandler{
		uc: uc,
	}
}

func SetRouteFunc(handler Handler) func(*gin.Engine) *gin.Engine {
	return func(engine *gin.Engine) *gin.Engine {
		routes := map[string]struct {
			fn     func(c *gin.Context)
			method string
		}{
			"/blogs":           {fn: fn(handler.ListBlogs), method: http.MethodGet},
			"/blogs/search":    {fn: fn(handler.SearchBlogs), method: http.MethodGet},
			"/comments/search": {fn: fn(handler.SearchComments), method: http.MethodGet},
		}

		api := engine.Group("/api")
		for path, route := range routes {
			api.Handle(route.method, path, route.fn)
		}

		return engine
	}
}

func fn(f func(ctx context.Context, c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if err := f(ctx, c); err != nil {
			handleError(ctx, c, err)
		}
	}
}

func handleError(ctx context.Context, c *gin.Context, err error) {
	log.Printf("failure: %v", err.Error())

	c.JSON(
		http.StatusInternalServerError,
		map[string]any{
			"error_message": err.Error(),
		},
	)
}

func (h *httpHandler) ListBlogs(ctx context.Context, c *gin.Context) error {
	res, err := h.uc.ListBlogs(ctx, c.Query("from"), c.Query("size"))
	if err != nil {
		return fmt.Errorf("failed to get blogs: %w", err)
	}

	c.JSON(http.StatusOK, res)

	return nil
}

func (h *httpHandler) SearchBlogs(ctx context.Context, c *gin.Context) error {
	res, err := h.uc.Search(ctx, "blogs", c.Query("query"), c.Query("from"), c.Query("size"))
	if err != nil {
		return fmt.Errorf("failed to search blogs: %w", err)
	}

	c.JSON(http.StatusOK, res)

	return nil
}

func (h *httpHandler) SearchComments(ctx context.Context, c *gin.Context) error {
	res, err := h.uc.Search(ctx, "comments", c.Query("query"), c.Query("from"), c.Query("size"))
	if err != nil {
		return fmt.Errorf("failed to search comments: %w", err)
	}

	c.JSON(http.StatusOK, res)

	return nil
}
