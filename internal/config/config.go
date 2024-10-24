package config

import (
	"github.com/rafalpienkowski/busgopher/internal/asb"
)

type Config struct {
	Connections map[string]asb.Connection `json:"connections"`
	Messages    map[string]asb.Message    `json:"messages"`
}

func Default() *Config {
	return &Config{
		Connections: make(map[string]asb.Connection),
		Messages:    make(map[string]asb.Message),
	}
}

type ConfigStorage interface {
	Load() (Config, error)
	Save(Config) error
}
