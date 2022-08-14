package order

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

var (
	ErrOrderInvalidAccrualResult = errors.New("got invalid accrual result")
	statusMapping                = map[AccrualStatus]ProcessingStatus{
		AccrualRegistered: New,
		AccrualProcessing: Processing,
		AccrualInvalid:    Invalid,
		AccrualProcessed:  Processed,
	}
	// _ app.OrderProcessor = (*Service)(nil)
)

type Repository interface {
	Create(ctx context.Context, ord Order) error
	ListByUser(ctx context.Context, usr user.User) ([]Order, error)
	ListUnprocessed(ctx context.Context) ([]Order, error)
	Update(ctx context.Context, ord Order) error
}

type Service struct {
	repo       Repository
	httpClient *resty.Client
	done       chan bool
}

func NewService(accrualAddr string, repo Repository) *Service {
	if repo == nil {
		panic("missing Repository, parameter must not be nil")
	}
	client := resty.New()
	client.SetBaseURL(accrualAddr).SetTimeout(time.Second).SetRetryCount(3).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusTooManyRequests
			},
		)
	s := &Service{httpClient: client, repo: repo}
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
	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := s.accrualUpdater(ctx)
			if err != nil {
				log.Err(err).Msg("caught error in accrualUpdaterService")
			}
		case <-s.done:
			return
		}
	}
}

func (s Service) accrualUpdater(ctx context.Context) error {
	log.Debug().Msg("orders updating...")
	orders, err := s.repo.ListUnprocessed(ctx)
	if err != nil {
		return err
	}

	errs, ctx := errgroup.WithContext(ctx)
	for _, order := range orders {
		ord := order
		errs.Go(func() error { return s.updateOrder(ctx, ord) })
	}
	return errs.Wait()
}

func (s Service) updateOrder(ctx context.Context, ord Order) error {
	accrual, err := s.getAccrual(ord)
	if err != nil {
		return err
	}
	if accrual.Status == AccrualProcessed {
		ord.Accrual = accrual.Accrual
	}
	ord.Status = statusMapping[accrual.Status]
	ord.Processed = time.Now()
	return s.repo.Update(ctx, ord)
}

func (s Service) getAccrual(ord Order) (Accrual, error) {
	var result Accrual
	_, err := s.httpClient.R().SetPathParams(map[string]string{
		"number": ord.Number.String(),
	}).SetResult(&result).Get("/api/orders/{number}")
	if err != nil {
		return result, err
	}
	log.Debug().Msgf("request accrual for order %s, got result (%v)",
		ord.Number.String(), result)
	if result.Order == "" {
		return result, ErrOrderInvalidAccrualResult
	}
	return result, nil
}
