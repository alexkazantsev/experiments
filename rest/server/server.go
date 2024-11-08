package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexkazantsev/experiments/rest/internal/user"
	"github.com/alexkazantsev/experiments/rest/server/middlewares"
)

type Middleware func(http.Handler) http.HandlerFunc

type Server struct {
	Addr   string
	Server *http.Server
}

func NewServer() *Server {
	var (
		addr     = ":8080"
		router   = NewRouter()
		userSrv  = user.NewUserService()
		userCtrl = user.NewUserController(userSrv)
		s        = &http.Server{
			Addr:    addr,
			Handler: router.V1,
		}
	)

	user.RegisterRoutes(router.V1, userCtrl)

	return &Server{
		Addr:   addr,
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

func Run() error {
	var (
		s    = NewServer()
		done = make(chan os.Signal, 1)
	)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	s.Use(
		middlewares.Recover,
		middlewares.Logger,
		middlewares.ContentType,
		middlewares.Timeout(2*time.Second),
		middlewares.User,
	)

	go func() {
		if err := s.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting the server: %+v\n", err)
		}

		log.Println("Stopped serving new connections.")
	}()

	log.Printf("server started at %s\n", s.Addr)

	<-done

	log.Println("Received signal to stop")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to stop: %+v\n", err)
	}

	log.Println("Server was gracefully stopped.")

	return nil
}
