package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/app"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/conf"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/auth"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/handler"
	midware "github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type Server struct {
	mart   *app.GopherMart
	srv    *http.Server
	router *chi.Mux
}

func NewServer(cfg conf.Server) (srv *Server, err error) {
	s := &Server{}
	s.mart = app.NewGopherMart(auth.NewServiceWithDefaultCredMan())

	s.router = chi.NewRouter()
	s.registerMiddlewares()
	s.registerHandlers()
	s.srv = &http.Server{
		Addr:    cfg.RunAddress,
		Handler: s.router,
	}
	return s, nil
}

func (s *Server) registerHandlers() {
	hApp := handler.NewApp(s.mart)
	// s.router.Route("/api", func(r chi.Router) {
	// 	s.router.Route("/user", func(r chi.Router) {
	// 		s.router.Post("/register", app.Service.RegisterUser)
	// 	})
	// })
	s.router.Post("/api/user/register", hApp.Auth.RegisterUser)
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
		// errors := repo.Close()
		// if errors != nil {
		// 	log.Error().Msgf("Caught an error due closing repository:%+v", errors)
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