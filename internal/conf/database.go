package conf

import (
	"errors"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	databaseFlag = "database-uri"
)

var ErrConfigDatabaseURINotSet = errors.New("connection string for DB is not set")

type Database struct {
	URI string
}

func (db *Database) SetPFlag() {
	pflag.StringP(databaseFlag, "d", "", "sets connection string for DB")
}

func (db *Database) Read() error {
	db.URI = viper.GetString(databaseFlag)
	if db.URI == "" {
		return ErrConfigDatabaseURINotSet
	}
	return nil
}
