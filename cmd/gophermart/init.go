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
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	cfg = conf.NewAppConfig()
	cfg.SetPFlag()

	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatal().Err(err).Interface("flag", pflag.CommandLine).Msg("can't bind argument")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	err = cfg.Read()
	if err != nil {
		log.Fatal().Err(err).Msg("can't read config")
	}
	log.Info().Str("address", cfg.RunAddress).Msg("cfg: server addr is set")
	log.Info().Msg("cfg: database uri is set")
	log.Info().Str("address", cfg.AccrualSystemAddress).Msg("cfg: accrual system addr is set")
}
