package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
)

// Config contains configurable values for the migration mechanism.
//type Config struct {
//	// DB config
//	Database string `env:"DB_CONNECTION_STRING" validate:"required"`
//	MigrationsPath   string `env:"DB_MIGRATIONS_PATH" validate:"required"`
//
//	// Server config
//	ConsoleServerAddress string `env:"CONSOLE_SERVER_ADDRESS" validate:"required"`
//}

// ReadConfig reads & validates config
func ReadConfig(v interface{}) error {
	err := env.ParseWithFuncs(v, nil)
	if err != nil {
		return err
	}

	validate := validator.New()

	err = validate.Struct(v)
	if err != nil {
		return err
	}

	return nil
}
