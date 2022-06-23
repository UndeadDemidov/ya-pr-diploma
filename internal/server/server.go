package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/auth"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/conf"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/handler"
	midware "github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type Server struct {
	srv    *http.Server
	router *chi.Mux
}

func NewServer(cfg conf.Server) (srv *Server, err error) {
	s := &Server{}
	s.router = chi.NewRouter()
	s.registerMiddlewares()
	s.registerHandlers()
	s.srv = &http.Server{Addr: cfg.RunAddress, Handler: s.router}
	return s, nil
}

func (s *Server) registerHandlers() {
	app := handler.NewApp(auth.Credentials{})
	// s.router.Route("/api", func(r chi.Router) {
	// 	s.router.Route("/user", func(r chi.Router) {
	// 		s.router.Post("/register", app.Auth.RegisterUser)
	// 	})
	// })
	s.router.Post("/api/user/register", app.Auth.RegisterUser)
}

func (s *Server) registerMiddlewares() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Compress(5))
	// r.Use(midware.Decompress)
	s.router.Use(midware.UserCookie)
}

func (s *Server) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("listen: %+v\n", err)
		}
	}()
	log.Info().Msg("Server started")

	<-ctx.Done()

	log.Info().Msg("Server stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// err := repo.Close()
		// if err != nil {
		// 	log.Error().Msgf("Caught an error due closing repository:%+v", err)
		// }

		log.Info().Msg("Everything is closed properly")
		cancel()
	}()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Error().Msgf("Server Shutdown Failed:%+v", err)
	}
	stop()
	log.Info().Msg("Server exited properly")
}
