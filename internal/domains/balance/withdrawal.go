package balance

import (
	"time"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/google/uuid"
)

type Withdrawal struct {
	ID        string            `json:"-"`
	User      user.User         `json:"-"`
	Order     primit.LuhnNumber `json:"order,string"`
	Sum       primit.Currency   `json:"sum"`
	Processed time.Time         `json:"processed_at"`
}

func NewWithdrawal(usr user.User, num primit.LuhnNumber, sum primit.Currency) (Withdrawal, error) {
	if !num.IsValid() {
		return Withdrawal{}, errors2.ErrOrderInvalidNumberFormat
	}
	return Withdrawal{
		ID:    uuid.New().String(),
		User:  usr,
		Order: num,
		Sum:   sum,
	}, nil
}
