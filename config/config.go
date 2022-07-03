package config

import (
	"encoding/json"
	"errors"
	"io"
)

var (
	ErrInvalidConfig = errors.New("Invalid Config")
)

type Config struct {
	Query  string
	ApiKey []string
}

func (c Config) isValid() bool {
	if len(c.ApiKey) < 1 {
		return false
	}
	if len(c.Query) < 1 {
		return false
	}
	return true
}

func (dc *Config) override(c Config) {
	if len(c.ApiKey) > 0 {
		dc.ApiKey = c.ApiKey
	}
	if len(c.Query) > 0 {
		dc.Query = c.Query
	}
}

func defaultConfig() Config {
	return Config{
		Query:  "vlog",
		ApiKey: []string{"<your api key>"},
	}
}

func New(r io.Reader) (c Config, err error) {
	dc := defaultConfig()

	err = json.NewDecoder(r).Decode(&c)
	if err != nil {
		return
	}
	if !c.isValid() {
		err = ErrInvalidConfig
		return
	}

	dc.override(c)

	return dc, nil
}
