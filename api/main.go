package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opensearch-project/opensearch-go/v2"

	"github.com/kqns91/searcher/api/handler"
	search "github.com/kqns91/searcher/api/repository/search"
	"github.com/kqns91/searcher/api/usecase"
)

var exitCode = 0

func main() {
	defer func() {
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	opensearchClient, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{os.Getenv("OPEN_SEARCH_ADDRESS")},
		Username:  os.Getenv("USER_NAME"),
		Password:  os.Getenv("PASSWORD"),
	})
	if err != nil {
		log.Printf("failed to create opensearch client: %v", err.Error())

		exitCode = 1

		return
	}

	searchRepo := search.New(opensearchClient)
	uc := usecase.New(searchRepo)
	h := handler.New(uc)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine = (handler.SetRouteFunc(h))(engine)

	log.Printf("listening and serving on port %s", os.Getenv("PORT"))

	if engine.Run(":" + os.Getenv("PORT")); err != nil {
		log.Printf("failed to serve: %v", err)

		exitCode = 1

		return
	}
}
