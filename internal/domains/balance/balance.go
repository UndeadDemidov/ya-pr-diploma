package balance

import (
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
)

type Balance struct {
	User      user.User       `json:"-"`
	Current   primit.Currency `json:"current"`
	Collected primit.Currency `json:"-"`
	Withdrawn primit.Currency `json:"withdrawn,omitempty"`
}

func NewBalance(usr user.User, cur, col, wth int64) Balance {
	return Balance{
		User:      usr,
		Current:   primit.Currency(cur),
		Collected: primit.Currency(col),
		Withdrawn: primit.Currency(wth),
	}
}
