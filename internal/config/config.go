package config

import (
	validate "github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

//Config...
type Config struct {
	PORT    string `validate:"required" envconfig:"PORT"`
	RUNTIME string `validate:"required" envconfig:"RUNTIME"`
	DB      *DB
}

// DB configuration
type DB struct {
	DNS string `validate:"required" envconfig:"DATABASE_URL"`
}

// LoadConf loads conf variables
func LoadConf(prefix string) (*Config, error) {

	c := new(Config)

	if err := envconfig.Process(prefix, c); err != nil {
		return nil, err
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

// validate configuration
func (c *Config) validate() error {

	var validator = validate.New()

	if err := validator.Struct(c); err != nil {
		return err
	}
	return nil
}
