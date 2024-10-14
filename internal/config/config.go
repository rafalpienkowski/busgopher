package config

import(
	"github.com/rafalpienkowski/busgopher/internal/asb"
)

type Config struct {
	Connections []asb.Connection `json:"connections"`
	Messages    []asb.Message    `json:"messages"`
}

func (config Config) Default() *Config {
	messages := []asb.Message{}
	connections := []asb.Connection{}

	return &Config{
		Connections: connections,
		Messages:    messages,
	}
}

type ConfigStorage interface {
	Load() (Config, error)
	Save(Config) error
}
