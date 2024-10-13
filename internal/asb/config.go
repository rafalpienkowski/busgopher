package asb

type Config struct {
	Connections *[]Connection `json:"connections"`
	Messages    *[]Message    `json:"messages"`
}

func (config Config) Default() *Config {
	messages := []Message{}
	connections := []Connection{}

	return &Config{
		Connections: &connections,
		Messages:    &messages,
	}
}

type ConfigStorage interface {
	Load() (*Config, error)
	Save(*Config) error
}
