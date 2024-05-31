package chocorpc

import (
	"context"
	"github.com/go-chocolate/chocolate/pkg/chocolate/cluster"
	"github.com/go-chocolate/chocolate/pkg/chocolate/cluster/endpoint"
	"net"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"github.com/go-chocolate/chocolate/pkg/chocolate/chocorpc/internal/interceptor"
)

type Server struct {
	config Config

	listener net.Listener

	services []func(*grpc.Server)

	options            []grpc.ServerOption
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	cluster.Cluster
}

func NewServer(config Config) *Server {
	var options = []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}
	options = append(options, config.apply()...)
	return &Server{
		options:           options,
		unaryInterceptors: []grpc.UnaryServerInterceptor{interceptor.Logger},
		config:            config,
	}
}

func (s *Server) Register(register func(*grpc.Server)) {
	s.services = append(s.services, register)
}

func (s *Server) WithServerOption(options ...grpc.ServerOption) *Server {
	s.options = append(s.options, options...)
	return s
}

func (s *Server) WithUnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor) *Server {
	s.unaryInterceptors = append(s.unaryInterceptors, interceptors...)
	return s
}

func (s *Server) WithStreamInterceptor(interceptors ...grpc.StreamServerInterceptor) *Server {
	s.streamInterceptors = append(s.streamInterceptors, interceptors...)
	return s
}

func (s *Server) Run(ctx context.Context) error {
	s.options = append(s.options,
		grpc.ChainUnaryInterceptor(s.unaryInterceptors...),
		grpc.ChainStreamInterceptor(s.streamInterceptors...),
	)
	server := grpc.NewServer(s.options...)
	for _, service := range s.services {
		service(server)
	}
	var err error
	if s.listener, err = net.Listen("tcp", s.config.Addr); err != nil {
		return err
	}

	if err := s.ClusterRegister(ctx, s.config.Name, s.endpoint()); err != nil {
		return err
	}

	return server.Serve(s.listener)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.listener.Close()
}

func (s *Server) endpoint() *endpoint.Endpoint {
	host, p, err := net.SplitHostPort(s.config.Addr)
	if err != nil {
		host = "0.0.0.0"
		p = "80"
	}
	port, _ := strconv.Atoi(p)

	return &endpoint.Endpoint{
		Protocol: endpoint.GRPC,
		Host:     host,
		Port:     uint16(port),
		Healthy: endpoint.HealthyOption{
			Enable:   true,
			Protocol: endpoint.GRPC,
			TLS:      false,
		},
	}
}
