package config

type InMemoryConfigStorage struct {
	Config Config
}

func (storage *InMemoryConfigStorage) Load() (Config, error) {

	return storage.Config, nil
}

func (storage *InMemoryConfigStorage) Save(config Config) error {
	storage.Config = config

	return nil
}
