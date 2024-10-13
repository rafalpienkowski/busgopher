package asb

import "errors"

type InMemoryConfigStorage struct {
	Config *Config
}

func (storage *InMemoryConfigStorage) Load() (*Config, error) {
	if storage.Config == nil {
        storage.Config = Config{}.Default()
	}

	return storage.Config, nil
}

func (storage *InMemoryConfigStorage) Save(config *Config) error {

    if config == nil {
        return errors.New("You're saving a nil config")
    }

    storage.Config = config

	return nil
}
