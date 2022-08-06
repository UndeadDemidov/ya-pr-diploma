package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
)

type ProcessingStatus int

var (
	_ fmt.Stringer   = (*ProcessingStatus)(nil)
	_ json.Marshaler = (*ProcessingStatus)(nil)
)

const (
	New ProcessingStatus = iota
	Processing
	Invalid
	Processed
)

var statuses = [...]string{"NEW", "PROCESSING", "INVALID", "PROCESSED"}

func (s ProcessingStatus) String() string {
	if s < New || s > Processed {
		return fmt.Sprintf("ProcessingStatus(%d)", int(s))
	}
	return statuses[s]
}

func (s ProcessingStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}

func (s ProcessingStatus) IsValid() bool {
	switch s {
	case New, Processing, Invalid, Processed:
		return true
	}
	return false
}

// Order
// Вообще-то это по смыслу не фига не заказ, а бонус за заказ! А баланс - это совокупность бонусов и списаний.
type Order struct {
	ID        string            `json:"-"`
	User      user.User         `json:"-"`
	Number    primit.LuhnNumber `json:"number,string"`
	Status    ProcessingStatus  `json:"status,string"`
	Accrual   primit.Currency   `json:"accrual,omitempty"`
	Unloaded  time.Time         `json:"uploaded_at"`
	Processed time.Time         `json:"-"`
}

func (o *Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	return json.Marshal(&struct {
		Unloaded string `json:"uploaded_at"`
		*Alias
	}{
		Unloaded: o.Unloaded.Format(time.RFC3339),
		Alias:    (*Alias)(o),
	})
}
