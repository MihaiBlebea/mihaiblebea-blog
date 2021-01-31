package api

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/MihaiBlebea/blog/go-broadcast/leads"
	"github.com/MihaiBlebea/blog/go-broadcast/page"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Private endpoint paths
var private = []string{
	"/lead",
}

// Custom errors
var (
	ErrPrivatePath = errors.New("Path is private")
)

type httpServer struct {
	pageService page.Service
	leadService leads.Service
	handler     http.Handler
	logger      *logrus.Logger
}

// HTTPServer interface
type HTTPServer interface {
	GetHandler() http.Handler
	TemplateHandler(w http.ResponseWriter, r *http.Request)
}

// NewHTTPServer returns a new http server service
func NewHTTPServer(
	pageService page.Service,
	leadService leads.Service,
	logger *logrus.Logger) HTTPServer {

	httpServer := httpServer{
		pageService: pageService,
		leadService: leadService,
		logger:      logger,
	}

	r := mux.NewRouter()

	r.Methods("POST").Path("/lead").HandlerFunc(httpServer.PostLeadHandler)

	r.PathPrefix("/static/").Handler(
		http.StripPrefix(
			"/static/",
			http.FileServer(
				http.Dir(
					httpServer.staticFolderPath(),
				),
			),
		),
	)

	r.PathPrefix("/").HandlerFunc(httpServer.TemplateHandler)

	httpServer.handler = r

	return &httpServer
}

func (h *httpServer) TemplateHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	reqID := h.requestID(8)

	h.logger.WithFields(logrus.Fields{
		"Url": path,
		"Id":  reqID,
	}).Info("Request started")
	start := time.Now()

	defer h.logger.WithFields(logrus.Fields{
		"Url":      path,
		"Id":       reqID,
		"Duration": time.Since(start),
	}).Info("Request ended")

	// Check if the path requested is not on the private list
	for _, p := range private {
		if p == path {
			h.logger.Error(ErrPrivatePath)
			page, _ := h.pageService.LoadErrorPage(ErrPrivatePath)

			err := page.Render(w)
			if err != nil {
				h.logger.Error(err)
			}

			return
		}
	}

	page, err := h.pageService.LoadTemplate(path)
	if err != nil {
		h.logger.Error(err)
	}

	err = page.Render(w)
	if err != nil {
		h.logger.Error(err)
	}
}

func (h *httpServer) PostLeadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := r.Form.Get("email")

	path := r.URL.Path
	reqID := h.requestID(8)

	h.logger.WithFields(logrus.Fields{
		"Url":   path,
		"Id":    reqID,
		"Email": email,
	}).Info("Request started")
	start := time.Now()

	defer h.logger.WithFields(logrus.Fields{
		"Url":      path,
		"Id":       reqID,
		"Duration": time.Since(start),
	}).Info("Request ended")

	err := h.leadService.Store(email)
	if err != nil {
		h.logger.Error(err)
	}

	page, err := h.pageService.LoadTemplate(path)
	if err != nil {
		h.logger.Error(err)
	}

	err = page.Render(w)
	if err != nil {
		h.logger.Error(err)
	}
}

func (h *httpServer) GetHandler() http.Handler {
	return h.handler
}

func (h *httpServer) staticFolderPath() string {
	p, err := os.Executable()
	if err != nil {
		h.logger.Fatal(err)
	}

	absPath := fmt.Sprintf(
		"%s/%s/",
		path.Dir(p),
		"static",
	)
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		h.logger.Fatal(err)
	}

	return absPath
}

func (h *httpServer) requestID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[random.Intn(len(charset))]
	}

	return string(b)
}
