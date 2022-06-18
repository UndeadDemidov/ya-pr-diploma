package conf

import (
	"errors"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	accrualSystemFlag = "accrual-system-address"
)

var ErrConfigAccrualSysAddrNotSet = errors.New("accrual system address is not set")
var _ Configurer = (*Externals)(nil)

type Externals struct {
	AccrualSystemAddress string
}

func (e *Externals) SetPFlag() {
	pflag.StringP(accrualSystemFlag, "r", "", "sets accrual system address")
}

func (e *Externals) Read() error {
	e.AccrualSystemAddress = viper.GetString(accrualSystemFlag)
	if e.AccrualSystemAddress == "" {
		return ErrConfigAccrualSysAddrNotSet
	}
	return nil
}
