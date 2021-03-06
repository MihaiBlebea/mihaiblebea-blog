package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MihaiBlebea/blog/go-broadcast/api"
	"github.com/MihaiBlebea/blog/go-broadcast/cache"
	"github.com/MihaiBlebea/blog/go-broadcast/leads"
	"github.com/MihaiBlebea/blog/go-broadcast/page"
	"github.com/MihaiBlebea/blog/go-broadcast/post"
	"github.com/sirupsen/logrus"
)

//go:generate go-bindata -o=assets/bindata.go --pkg=assets static/templates/... static/markdown/...

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	cache := cache.New()
	postService := post.New(
		isDev(),
	)
	pageService := page.New(postService, cache, logger)

	leadService := leads.New(
		os.Getenv("MAILCHIMP_API_KEY"),
		os.Getenv("MAILCHIMP_LIST_ID"),
	)

	server := api.NewHTTPServer(pageService, leadService, logger)

	httpPort := fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))
	logger.Info(fmt.Sprintf("Application started HTTP server on port %s", httpPort))

	err := http.ListenAndServe(httpPort, server.GetHandler())
	if err != nil {
		logger.Fatal(err)
	}
}

func isDev() bool {
	if os.Getenv("DEV") == "true" {
		return true
	}

	return false
}
