package balance

import (
	"time"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/order"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
)

type Balance struct {
	User      user.User       `json:"-"`
	Current   primit.Currency `json:"current"`
	Collected primit.Currency `json:"-"`
	Withdrawn primit.Currency `json:"withdrawn,omitempty"`
}

type Withdrawal struct {
	User      user.User       `json:"-"`
	Order     order.Order     `json:"order"` // custom marshaler
	Sum       primit.Currency `json:"sum"`
	Processed time.Time       `json:"processed_at"` // custom marshaler
}
