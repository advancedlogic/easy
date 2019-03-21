package transport

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
	"github.com/zsais/go-gin-prometheus"
	"net/http"
	"time"
)

type Option func(*Rest) error

const (
	GET    = "get"
	POST   = "post"
	PUT    = "put"
	DELETE = "delete"
)

func WithPort(port int) Option {
	return func(rest *Rest) error {
		if port > 0 {
			rest.port = port
			return nil
		}
		return errors.New("port cannot be zero")
	}
}

func WithHandler(mode, route string, handler interface{}) Option {
	return func(rest *Rest) error {
		return rest.Handler(mode, route, handler)
	}
}

func WithMiddleware(middleware interface{}) Option {
	return func(rest *Rest) error {
		return rest.Middleware(middleware)
	}
}
func WithStaticFilesFolder(route, folder string) Option {
	return func(rest *Rest) error {
		return rest.StaticFilesFolder(route, folder)
	}
}

// Rest Server
type Rest struct {
	// Port bound to server
	port           int
	readTimeout    time.Duration
	writeTimeout   time.Duration
	getHandlers    map[string][]gin.HandlerFunc
	postHandlers   map[string][]gin.HandlerFunc
	putHandlers    map[string][]gin.HandlerFunc
	deleteHandlers map[string][]gin.HandlerFunc
	middleware     []func(ctx *gin.Context)
	websiteFolder  map[string]string
	cert           string
	key            string
	server         *http.Server
	logger         *logrus.Logger
}

func NewRest(options ...Option) (*Rest, error) {
	rest := &Rest{
		port:           8080,
		getHandlers:    make(map[string][]gin.HandlerFunc),
		postHandlers:   make(map[string][]gin.HandlerFunc),
		putHandlers:    make(map[string][]gin.HandlerFunc),
		deleteHandlers: make(map[string][]gin.HandlerFunc),
		middleware:     make([]func(ctx *gin.Context), 0),
		websiteFolder:  make(map[string]string),
		logger:         logrus.New(),
	}

	for _, option := range options {
		err := option(rest)
		if err != nil {
			return nil, err
		}
	}

	return rest, nil
}
func (r *Rest) Handler(mode, route string, handler interface{}) error {
	switch mode {
	case GET:
		if _, exists := r.getHandlers[route]; !exists {
			r.getHandlers[route] = make([]gin.HandlerFunc, 0)
		}
		r.getHandlers[route] = append(r.getHandlers[route], gin.HandlerFunc(handler.(func(*gin.Context))))
	case POST:
		if _, exists := r.postHandlers[route]; !exists {
			r.postHandlers[route] = make([]gin.HandlerFunc, 0)
		}
		r.postHandlers[route] = append(r.postHandlers[route], gin.HandlerFunc(handler.(func(*gin.Context))))
	case PUT:
		if _, exists := r.putHandlers[route]; !exists {
			r.putHandlers[route] = make([]gin.HandlerFunc, 0)
		}
		r.putHandlers[route] = append(r.putHandlers[route], gin.HandlerFunc(handler.(func(*gin.Context))))
	case DELETE:
		if _, exists := r.deleteHandlers[route]; !exists {
			r.deleteHandlers[route] = make([]gin.HandlerFunc, 0)
		}
		r.deleteHandlers[route] = append(r.deleteHandlers[route], gin.HandlerFunc(handler.(func(*gin.Context))))
	default:
		return errors.New("possible modes are GET, POST, PUT, DELETE")
	}
	return nil
}

func (r *Rest) Middleware(middleware interface{}) error {
	if r.middleware == nil {
		r.middleware = make([]func(*gin.Context), 0)
	}
	r.middleware = append(r.middleware, middleware.(func(ctx *gin.Context)))
	return nil
}

func (r *Rest) StaticFilesFolder(uri, folder string) error {
	r.websiteFolder[uri] = folder
	return nil
}

func (r *Rest) Run() error {
	router := gin.New()
	router.Use(ginlogrus.Logger(r.logger), gin.Recovery())
	router.GET("/healthcheck", func(c *gin.Context) {
		c.String(200, "product service is good")
	})

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	for _, middleware := range r.middleware {
		router.Use(gin.HandlerFunc(middleware))
	}

	for path, handlers := range r.getHandlers {
		router.GET(path, handlers...)
	}

	for path, handler := range r.postHandlers {
		router.POST(path, handler...)
	}
	for path, handler := range r.putHandlers {
		router.PUT(path, handler...)
	}

	for path, handler := range r.deleteHandlers {
		router.DELETE(path, handler...)
	}

	if len(r.websiteFolder) > 0 {
		for uri, folder := range r.websiteFolder {
			router.Static(uri, folder)
		}
	}

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", r.port),
		Handler:        router,
		ReadTimeout:    r.readTimeout,
		WriteTimeout:   r.writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	r.server = s
	go func() {
		if r.cert != "" && r.key != "" {
			if err := s.ListenAndServeTLS(r.cert, r.key); err != nil {
				r.logger.Fatal(err)
			}

		} else if err := s.ListenAndServe(); err != nil {
			r.logger.Fatal(err)
		}
	}()
	return nil
}

func (r *Rest) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.server.Shutdown(ctx)
}
