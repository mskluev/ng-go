package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"example.com/user/ng-go/internal/route"
)

// Handler serves various HTTP endpoints of the Prometheus server
type Handler struct {
	logger  log.Logger
	router  *route.Router
	options *Options
	birth   time.Time
	cwd     string
}

// Options for the web Handler.
type Options struct {
	Context context.Context
	Flags   map[string]string

	ListenAddress  string
	CORSOrigin     *regexp.Regexp
	ReadTimeout    time.Duration
	MaxConnections int
	ExternalURL    *url.URL
	RoutePrefix    string
	PageTitle      string
}

// New initializes a new web Handler.
func New(logger log.Logger, o *Options) *Handler {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	router := route.New().
		WithInstrumentation(setPathWithPrefix(""))

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "<error retrieving current working directory>"
	}

	h := &Handler{
		logger:  logger,
		router:  router,
		options: o,
		birth:   time.Now().UTC(),
		cwd:     cwd,
	}

	router.Get("/status", h.status)

	return h
}

// Run serves the HTTP endpoints
func (h *Handler) Run() {
	level.Info(h.logger).Log("msg", "Start listening for connections", "address", h.options.ListenAddress)
	http.ListenAndServe(h.options.ListenAddress, h.router)
}

func (h *Handler) status(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Birth          time.Time
		CWD            string
		GoroutineCount int
		GOMAXPROCS     int
		GOGC           string
		GODEBUG        string
	}{
		Birth:          h.birth,
		CWD:            h.cwd,
		GoroutineCount: runtime.NumGoroutine(),
		GOMAXPROCS:     runtime.GOMAXPROCS(0),
		GOGC:           os.Getenv("GOGC"),
		GODEBUG:        os.Getenv("GODEBUG"),
	}

	dec := json.NewEncoder(w)
	if err := dec.Encode(status); err != nil {
		http.Error(w, fmt.Sprintf("error encoding JSON: %s", err), http.StatusInternalServerError)
	}
}

func setPathWithPrefix(prefix string) func(handlerName string, handler http.HandlerFunc) http.HandlerFunc {
	return func(handlerName string, handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			handler(w, r.WithContext(ContextWithPath(r.Context(), prefix+r.URL.Path)))
		}
	}
}

type pathParam struct{}

// ContextWithPath returns a new context with the given path to be used later
// when logging the query.
func ContextWithPath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, pathParam{}, path)
}
