package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/app"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/conf"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/auth"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/infra/persist/postgre"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/handler"
	midware "github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type Server struct {
	dbPool *pgxpool.Pool
	mart   *app.GopherMart
	srv    *http.Server
	router *chi.Mux
}

func NewServer(cfg *conf.App) (srv *Server, err error) {
	if cfg == nil {
		return nil, errors.New("conf.App is nil")
	}
	s := &Server{}
	ctx := context.Background()
	s.dbPool, err = pgxpool.Connect(ctx, cfg.Database.URI)
	if err != nil {
		return nil, err
	}

	repo, err := postgre.NewPersist(ctx, s.dbPool)
	if err != nil {
		return nil, err
	}
	// ToDo конфигуратор?
	a := auth.NewServiceWithDefaultCredMan(repo.Auth, user.NewService(repo.User))
	s.mart = app.NewGopherMart(a)

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
		s.dbPool.Close()
		log.Info().Msg("Everything is closed properly")
		cancel()
	}()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Error().Msgf("Server Shutdown Failed:%+v", err)
	}
	stop()
	log.Info().Msg("Server exited properly")
}
