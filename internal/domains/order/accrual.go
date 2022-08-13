package order

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
)

type AccrualStatus int

const (
	AccrualRegistered AccrualStatus = iota
	AccrualProcessing
	AccrualInvalid
	AccrualProcessed
)

var (
	_ fmt.Stringer     = (*AccrualStatus)(nil)
	_ json.Unmarshaler = (*AccrualStatus)(nil)

	ErrOrderInvalidAccrualStatus = errors.New("invalid order processing status")
)

func ParseAccrualStatus(s string) (AccrualStatus, error) {
	var strings = map[string]AccrualStatus{
		"REGISTERED": AccrualRegistered,
		"PROCESSING": AccrualProcessing,
		"INVALID":    AccrualInvalid,
		"PROCESSED":  AccrualProcessed,
	}
	if status, ok := strings[s]; ok {
		return status, nil
	}
	return AccrualInvalid, ErrOrderInvalidAccrualStatus
}

func (a AccrualStatus) String() string {
	if a < AccrualRegistered || a > AccrualProcessed {
		return fmt.Sprintf("AccrualStatus(%d)", int(a))
	}
	var statuses = [...]string{"REGISTERED", "PROCESSING", "INVALID", "PROCESSED"}
	return statuses[a]
}

func (a *AccrualStatus) UnmarshalJSON(data []byte) (err error) {
	var v string
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	if *a, err = ParseAccrualStatus(v); err != nil {
		return err
	}
	return nil
}

type Accrual struct {
	Order   string          `json:"order"`
	Status  AccrualStatus   `json:"status"`
	Accrual primit.Currency `json:"accrual"`
}
