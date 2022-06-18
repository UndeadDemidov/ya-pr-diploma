package conf

import (
	"errors"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	runAddressFlag = "run-address"
)

var ErrConfigRunAddressNotSet = errors.New("server address is not set")
var _ Configurer = (*Server)(nil)

type Server struct {
	RunAddress string
}

func (s *Server) SetPFlag() {
	pflag.StringP(runAddressFlag, "a", ":8080", "sets http server address")
}

func (s *Server) Read() error {
	s.RunAddress = viper.GetString(runAddressFlag)
	if s.RunAddress == "" {
		return ErrConfigRunAddressNotSet
	}
	return nil
}
