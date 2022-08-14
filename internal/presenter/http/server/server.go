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
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/balance"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/order"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/infra/persist/postgre"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/handler"
	midware "github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type Closer interface {
	Close()
}

var _ Closer = (*app.GopherMart)(nil)

type Server struct {
	dbPool   *pgxpool.Pool
	mart     *app.GopherMart
	srv      *http.Server
	router   *chi.Mux
	sessions *midware.Sessions
}

func NewServer(cfg *conf.App) (srv *Server, err error) {
	if cfg == nil {
		panic("missing *conf.App, parameter must not be nil")
	}
	s := &Server{}

	// ToDo конфигуратор?
	ctx := context.Background()
	s.dbPool, err = pgxpool.Connect(ctx, cfg.Database.URI)
	if err != nil {
		return nil, err
	}

	repo, err := postgre.NewPersist(ctx, s.dbPool)
	if err != nil {
		return nil, err
	}
	svcAuth := auth.NewServiceWithDefaultCredMan(repo.Auth, user.NewService(repo.User))
	svcOrder := order.NewService(cfg.AccrualSystemAddress, repo.Order)
	svcBalance := balance.NewService(repo.Balance)
	// app configuration
	s.mart = app.NewGopherMart(svcAuth, svcOrder, svcBalance)
	// router configuration
	s.sessions = midware.NewDefaultSessions()
	s.router = s.buildRouter(
		handler.NewAuth(s.mart, s.sessions),
		handler.NewOrder(s.mart),
		handler.NewBalance(s.mart),
	)

	s.srv = &http.Server{
		Addr:    cfg.RunAddress,
		Handler: s.router,
	}
	return s, nil
}

func (s *Server) buildRouter(auth *handler.Auth, order *handler.Order, bal *handler.Balance) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(midware.Decompress)

	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", auth.RegisterUser)
		r.Post("/api/user/login", auth.LoginUser)
	})
	r.Group(func(r chi.Router) {
		r.Use(midware.SessionsCookie(s.sessions))
		r.Post("/api/user/orders", order.UploadOrder)
		r.Get("/api/user/orders", order.DownloadOrders)
		r.Get("/api/user/balance", bal.Get)
	})
	return r
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
		s.mart.Close()
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
