package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/advancedlogic/easy/commons"
	"github.com/advancedlogic/easy/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
	"github.com/zsais/go-gin-prometheus"
	"net"
	"net/http"
	"strings"
	"time"
)

func WithPort(port int) interfaces.TransportOption {
	return func(t interfaces.Transport) error {
		if port > 0 {
			rest := t.(*Rest)
			rest.port = port
			return nil
		}
		return errors.New("port cannot be zero")
	}
}

func WithHandler(mode, route string, handler interface{}) interfaces.TransportOption {
	return func(t interfaces.Transport) error {
		rest := t.(*Rest)
		return rest.Handler(mode, route, handler)
	}
}

func GET(route string, handler interface{}) interfaces.TransportOption {
	return func(t interfaces.Transport) error {
		rest := t.(*Rest)
		option := WithHandler(commons.ModeGet, route, handler)
		return option(rest)
	}
}

func POST(route string, handler interface{}) interfaces.TransportOption {
	return func(t interfaces.Transport) error {
		rest := t.(*Rest)
		option := WithHandler(commons.ModePost, route, handler)
		return option(rest)
	}
}

func PUT(route string, handler interface{}) interfaces.TransportOption {
	return func(t interfaces.Transport) error {
		rest := t.(*Rest)
		option := WithHandler(commons.ModePut, route, handler)
		return option(rest)
	}
}

func DELETE(route string, handler interface{}) interfaces.TransportOption {
	return func(t interfaces.Transport) error {
		rest := t.(*Rest)
		option := WithHandler(commons.ModeDelete, route, handler)
		return option(rest)
	}
}

func WithMiddleware(middleware interface{}) interfaces.TransportOption {
	return func(t interfaces.Transport) error {
		rest := t.(*Rest)
		return rest.Middleware(middleware)
	}
}
func WithStaticFilesFolder(route, folder string) interfaces.TransportOption {
	return func(t interfaces.Transport) error {
		rest := t.(*Rest)
		return rest.StaticFilesFolder(route, folder)
	}
}

func WithLogger(logger *logrus.Logger) interfaces.TransportOption {
	return func(t interfaces.Transport) error {
		rest := t.(*Rest)
		return rest.WithLogger(logger)
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
	router         *gin.Engine
	*logrus.Logger
}

func New(options ...interfaces.TransportOption) (*Rest, error) {
	rest := &Rest{
		port:           8080,
		getHandlers:    make(map[string][]gin.HandlerFunc),
		postHandlers:   make(map[string][]gin.HandlerFunc),
		putHandlers:    make(map[string][]gin.HandlerFunc),
		deleteHandlers: make(map[string][]gin.HandlerFunc),
		middleware:     make([]func(ctx *gin.Context), 0),
		websiteFolder:  make(map[string]string),
		router:         gin.New(),
		Logger:         logrus.New(),
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
	case commons.ModeGet:
		if _, exists := r.getHandlers[route]; !exists {
			r.getHandlers[route] = make([]gin.HandlerFunc, 0)
		}
		r.getHandlers[route] = append(r.getHandlers[route], gin.HandlerFunc(handler.(func(*gin.Context))))
	case commons.ModePost:
		if _, exists := r.postHandlers[route]; !exists {
			r.postHandlers[route] = make([]gin.HandlerFunc, 0)
		}
		r.postHandlers[route] = append(r.postHandlers[route], gin.HandlerFunc(handler.(func(*gin.Context))))
	case commons.ModePut:
		if _, exists := r.putHandlers[route]; !exists {
			r.putHandlers[route] = make([]gin.HandlerFunc, 0)
		}
		r.putHandlers[route] = append(r.putHandlers[route], gin.HandlerFunc(handler.(func(*gin.Context))))
	case commons.ModeDelete:
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
	router := r.router
	router.Use(ginlogrus.Logger(r.Logger), gin.Recovery())
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

	if err := r.findAlternativePort(); err != nil {
		r.Fatal(err)
	}

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", r.port),
		Handler:        router,
		ReadTimeout:    r.readTimeout,
		WriteTimeout:   r.writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	r.server = s
	httpHeader := "http"
	go func() {
		if r.cert != "" && r.key != "" {
			if err := s.ListenAndServeTLS(r.cert, r.key); err != nil {
				r.Fatal(err)
			}
			httpHeader += "s"
		} else if err := s.ListenAndServe(); err != nil {
			r.Fatal(err)
		}

	}()
	r.Info(fmt.Sprintf("Http(s) server listening on port %d", r.port))
	return nil
}

func (r *Rest) scanPort(ip string, port int, timeout time.Duration) error {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			err = r.scanPort(ip, port, timeout)
		} else {
			fmt.Println(port, "closed")
		}
		return err
	}

	if err = conn.Close(); err != nil {
		return err
	}
	fmt.Println(port, "open")
	return nil
}

func (r *Rest) findAlternativePort() error {
	currentPort := r.port
	for port := currentPort; port < 32000; port++ {
		err := r.scanPort("localhost", port, 10*time.Second)
		if err != nil {
			r.port = port
			return nil
		}
	}
	return errors.New("no alternatives port found")
}

func (r *Rest) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.server.Shutdown(ctx)
}

func (r *Rest) WithLogger(logger *logrus.Logger) error {
	if logger != nil {
		r.Logger = logger
		return nil
	}
	return errors.New("logger cannot be nil")
}

func (r *Rest) Router() (interface{}, error) {
	if r.router != nil {
		return r.router, nil
	}
	return nil, errors.New("router is nil")
}
