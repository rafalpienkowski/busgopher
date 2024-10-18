package config

import (
	"github.com/rafalpienkowski/busgopher/internal/asb"
)

type Config struct {
	NConnections map[string]asb.Connection `json:"connections"`
	NMessages    map[string]asb.Message    `json:"messages"`
}

func (config Config) Default() *Config {
	return &Config{
		NConnections: make(map[string]asb.Connection),
		NMessages:    make(map[string]asb.Message),
	}
}

type ConfigStorage interface {
	Load() (Config, error)
	Save(Config) error
}
