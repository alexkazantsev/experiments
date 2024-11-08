package server

import (
	"log"
	"net/http"
	"time"

	"github.com/alexkazantsev/experiments/rest/server/middlewares"
)

type Middleware func(http.Handler) http.HandlerFunc

type Server struct {
	Addr   string
	Router *http.ServeMux
	Server *http.Server
}

func NewServer(router *Router) *Server {
	var (
		addr = ":8080"
		s    = &http.Server{
			Addr:    addr,
			Handler: router.V1,
		}
	)

	return &Server{
		Addr:   ":8080",
		Router: router.V1,
		Server: s,
	}
}

func (s *Server) Use(m ...Middleware) {
	var (
		n = len(m)
		i = n - 1
	)

	for i >= 0 {
		s.Server.Handler = m[i](s.Server.Handler)
		i--
	}
}

func (s *Server) Shutdown() error {
	return nil
}

func Run(s *Server) error {
	s.Use(
		middlewares.Recover,
		middlewares.Logger,
		middlewares.ContentType,
		middlewares.Timeout(2*time.Second),
		middlewares.User,
	)

	log.Printf("server started at %s", s.Addr)

	return s.Server.ListenAndServe()

}
