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
	Query    string   `json:"query"`
	ApiKey   []string `json:"apiKey"`
	MongoUri string   `json:"mongoUri"` // mongo uri
	DBName   string   `json:"dbName"`
	HttpPort int      `json:"httpPort"`
}

func (c Config) isValid() bool {
	if len(c.ApiKey) < 1 {
		return false
	}
	if len(c.Query) < 1 {
		return false
	}
	if len(c.MongoUri) < 1 {
		return false
	}
	if len(c.DBName) < 1 {
		return false
	}
	if c.HttpPort < 1 {
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
	if len(c.MongoUri) > 0 {
		dc.MongoUri = c.MongoUri
	}
	if len(c.DBName) > 0 {
		dc.DBName = c.DBName
	}
	if c.HttpPort > 0 {
		dc.HttpPort = c.HttpPort
	}
}

func defaultConfig() Config {
	return Config{
		Query:    "vlog",
		ApiKey:   []string{"<your api key>"},
		DBName:   "fampay",
		HttpPort: 3000,
		MongoUri: "mongodb://localhost:27017",
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
