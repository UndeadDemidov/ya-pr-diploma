package order

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// var _ app.OrderProcessor = (*Service)(nil)

type Repository interface {
	Create(ctx context.Context, ord Order) error
	ListByUser(ctx context.Context, usr user.User) ([]Order, error)
	ListUnprocessed(ctx context.Context) ([]Order, error)
}

type Service struct {
	accrualSystemAddress string
	repo                 Repository
	done                 chan bool
}

func NewService(accrualAddr string, repo Repository) *Service {
	if repo == nil {
		panic("missing Repository, parameter must not be nil")
	}
	s := &Service{accrualSystemAddress: accrualAddr, repo: repo}
	s.done = make(chan bool)
	go s.accrualUpdaterService()
	return s
}

func (s Service) Add(ctx context.Context, usr user.User, num string) error {
	num64, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return err
	}
	lnum := primit.LuhnNumber(num64)
	ord, err := NewOrder(usr, lnum)
	if err != nil {
		return err
	}
	err = s.repo.Create(ctx, ord)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) List(ctx context.Context, usr user.User) (ords []Order, err error) {
	return s.repo.ListByUser(ctx, usr)
}

func (s Service) Close() {
	s.done <- true
}

func (s Service) accrualUpdaterService() {
	ctx := context.Background()
	tick := time.Tick(4 * time.Second)
	for {
		select {
		case <-tick:
			err := s.accuralUpdater(ctx)
			if err != nil {
				log.Err(err).Msg("caught error in accrualUpdaterService")
			}
		case <-s.done:
			return
		}
	}
}

func (s Service) accuralUpdater(ctx context.Context) error {
	log.Debug().Msg("orders updating...")
	orders, err := s.repo.ListUnprocessed(ctx)
	if err != nil {
		return err
	}

	errs, ctx := errgroup.WithContext(ctx)
	for _, order := range orders {
		errs.Go(func() error { return s.updateOrder(ctx, order) }) //nolint:govet
	}
	return errs.Wait()
}

func (s Service) updateOrder(ctx context.Context, ord Order) error {
	var result struct {
		Order   string          `json:"order"`
		Status  string          `json:"status"`
		Accrual primit.Currency `json:"accrual"`
	}
	client := resty.New()
	client.SetBaseURL(s.accrualSystemAddress).SetTimeout(time.Second).SetRetryCount(3).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusTooManyRequests
			},
		)
	_, err := client.R().SetPathParams(map[string]string{
		"number": ord.Number.String(),
	}).SetResult(&result).Get("/api/orders/{number}")
	if err != nil {
		return err
	}
	log.Debug().Msgf("request accrual for order %s, got result (%v)",
		ord.Number.String(), result)
	return nil
}
