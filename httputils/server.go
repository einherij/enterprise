package httputils

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const shutdownTimeout = 10 * time.Second

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type Server struct {
	name     string
	server   HTTPServer
	listener net.Listener
}

func NewServer(name, port string, handler http.Handler) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, fmt.Errorf("cannot listen server port: %w", err)
	}
	return &Server{
		name:     name,
		server:   &http.Server{Handler: handler},
		listener: ln,
	}, nil
}

func (s *Server) Run(ctx context.Context) {
	go func() {
		logrus.Infof("starting %v server on %s", s.name, s.listener.Addr().String())
		err := s.server.Serve(s.listener)
		if err != nil && err != http.ErrServerClosed {
			logrus.Errorf("error serving %v web: %+v", s.name, err)
		}
	}()

	<-ctx.Done()
	{
		// shutdown
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		err := s.server.Shutdown(ctx)
		if err != nil {
			logrus.Infof("error shutting down http server: %v", err)
		}
	}
}

func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}
