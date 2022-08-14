package main

import (
	"os"
	"strings"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/conf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfg *conf.App

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	cfg = conf.NewAppConfig()
	cfg.SetPFlag()

	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatal().Err(err).Msgf("can't bind argument flags %v", pflag.CommandLine)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	err = cfg.Read()
	if err != nil {
		log.Fatal().Err(err).Msg("can't read config")
	}
	log.Info().Msgf("cfg: server addr is set to %v", cfg.RunAddress)
	log.Info().Msgf("cfg: database uri is set to %v", cfg.URI)
	log.Info().Msgf("cfg: accrual system addr is set to %v", cfg.AccrualSystemAddress)
}
