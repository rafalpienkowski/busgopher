package config

import (
	"github.com/rafalpienkowski/busgopher/internal/asb"
)

type Config struct {
	Connections []asb.Connection `json:"connections"`
	Messages    []asb.Message    `json:"messages"`

	NConnections map[string]asb.Connection `json:"nconnections"`
	NMessages    map[string]asb.Message    `json:"nmessages"`
}

func (config Config) Default() *Config {
	messages := []asb.Message{}
	connections := []asb.Connection{}

	return &Config{
		Connections:  connections,
		Messages:     messages,
		NConnections: make(map[string]asb.Connection),
		NMessages:    make(map[string]asb.Message),
	}
}

type ConfigStorage interface {
	Load() (Config, error)
	Save(Config) error
}
