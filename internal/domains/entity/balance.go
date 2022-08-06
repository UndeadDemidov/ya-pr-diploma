package entity

import (
	"encoding/json"
	"time"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
)

type Balance struct {
	User      user.User       `json:"-"`
	Current   primit.Currency `json:"current"`
	Collected primit.Currency `json:"-"`
	Withdrawn primit.Currency `json:"withdrawn,omitempty"`
}

var _ json.Marshaler = (*Withdrawal)(nil)

type Withdrawal struct {
	ID        string          `json:"-"`
	User      user.User       `json:"-"`
	Order     Order           `json:"order"`
	Sum       primit.Currency `json:"sum"`
	Processed time.Time       `json:"processed_at"`
}

func (w *Withdrawal) MarshalJSON() ([]byte, error) {
	type Alias Withdrawal
	return json.Marshal(&struct {
		Order string `json:"order"`
		*Alias
	}{
		Order: w.Order.Number.String(),
		Alias: (*Alias)(w),
	})
}
