package primit

import (
	"encoding/json"
	"fmt"
)

type Currency int64

var (
	_ fmt.Stringer   = (*Currency)(nil)
	_ json.Marshaler = (*Currency)(nil)
)

func Float64ToCurrency(f float64) Currency {
	return Currency((f * 100) + 0.5)
}

func (c Currency) Float64() float64 {
	return float64(c) / 100
}

func (c Currency) String() string {
	if c == c/100*100 {
		return fmt.Sprintf("%.0f", c.Float64())
	}
	return fmt.Sprintf("%.2f", c.Float64())
}

func (c Currency) MarshalJSON() ([]byte, error) {
	return []byte(c.String()), nil
}

func (c *Currency) UnmarshalJSON(data []byte) error {
	var v float64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*c = Float64ToCurrency(v)
	return nil
}
