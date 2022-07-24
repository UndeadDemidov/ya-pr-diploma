package order

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
)

type ProcessingStatus int

var _ fmt.Stringer = (*ProcessingStatus)(nil)

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

type Order struct {
	ID        string           `json:"-"`
	User      user.User        `json:"-"`
	Number    LuhnNumber       `json:"number,string"`
	Status    ProcessingStatus `json:"status,string"`
	Accrual   int64            `json:"accrual,omitempty"`
	Unloaded  time.Time        `json:"uploaded_at"`
	Processed time.Time        `json:"-"`
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
