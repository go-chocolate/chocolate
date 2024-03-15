package chocohttp

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	"github.com/go-chocolate/chocolate/pkg/chocolate/chocohttp/internal/handler"
	"github.com/go-chocolate/chocolate/pkg/chocolate/chocohttp/internal/middleware"
)

type Server struct {
	srv     *http.Server
	handler http.Handler
	config  Config
}

func NewServer(config Config) *Server {
	if config.Addr == "" {
		config.Addr = ":8080"
	}
	return &Server{
		config: config,
		srv:    &http.Server{Addr: config.Addr},
	}
}

func (s *Server) Run(ctx context.Context) error {
	middlewares := s.middlewares()
	for _, middle := range middlewares {
		s.handler = middle(s.handler)
	}
	s.srv.Handler = s.handler
	if s.config.TLS != nil {
		return s.srv.ListenAndServeTLS(s.config.TLS.CertFile, s.config.TLS.KeyFile)
	} else {
		return s.srv.ListenAndServe()
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) Router(router func(mux *http.ServeMux)) {
	s.check()
	mux := http.NewServeMux()
	mux.HandleFunc(handler.HealthPath, handler.Health)
	router(mux)
	s.handler = mux
}

func (s *Server) Httprouter(router func(router *httprouter.Router)) {
	s.check()
	e := httprouter.New()
	e.HandlerFunc(http.MethodGet, handler.HealthPath, handler.Health)
	router(e)
	s.handler = e
}

func (s *Server) GINRouter(router func(router gin.IRouter)) {
	s.check()
	e := gin.New()
	e.GET(handler.HealthPath, func(c *gin.Context) { handler.Health(c.Writer, c.Request) })
	router(e)
	s.handler = e
}

func (s *Server) ListenOn() string {
	return s.srv.Addr
}

func (s *Server) check() {
	if s.handler != nil {
		t := reflect.TypeOf(s.handler)
		logrus.Panic(fmt.Sprintf("http hander has been setted with %v", t))
	}
}

func (s *Server) middlewares() []middleware.Middleware {
	var middlewares = []middleware.Middleware{
		middleware.Recovery(),
		middleware.TraceId(),
		middleware.Trace(s.config.Name),
	}
	if s.config.Options.Logger.Enable {
		middlewares = append(middlewares, middleware.Logger())
	}
	if s.config.Options.CORS.Enable {
		middlewares = append(middlewares, middleware.CORS(s.config.Options.CORS.build()))
	}
	if s.config.Options.RateLimit.Enable {
		middlewares = append(middlewares, middleware.RateLimit(s.config.Options.RateLimit.Limit))
	}
	return middlewares
}
