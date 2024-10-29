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

func GetTestConfig() Config {
	connections := make(map[string]asb.Connection)
	connections["test-connection"] = asb.Connection{
		Namespace: "test.azure.com",
		Destinations: []string{
			"queue",
			"topic",
		},
	}
	messages := make(map[string]asb.Message)
	messages["test-message"] = asb.Message{
		Body: "{ test msg body }",
	}

	return Config{
		Connections: connections,
		Messages:    messages,
	}
}

type ConfigStorage interface {
	Load() (Config, error)
	Save(Config) error
}
