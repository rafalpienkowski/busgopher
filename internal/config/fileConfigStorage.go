package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/rafalpienkowski/busgopher/internal/asb"
)

const configName = "config.json"

type FileConfigStorage struct{}

func (storage *FileConfigStorage) Load() (Config, error) {
	var config *Config

	bytes, err := readFile(configName)
	if err != nil {
		return Config{}, err
	}

	if len(bytes) == 0 {
		json, jerr := json.Marshal(config.Default())
		if jerr != nil {
			return Config{}, jerr
		}

		err = writeFile(configName, string(json))
		if err != nil {
			return Config{}, err
		}

		bytes = []byte(json)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return Config{
			Messages:    make(map[string]asb.Message),
			Connections: make(map[string]asb.Connection),
		}, err
	}

	return *config, nil
}

func (storage *FileConfigStorage) Save(config Config) error {
	json, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = writeFile(configName, string(json))
	if err != nil {
		return err
	}
	return nil
}

func readFile(filePath string) ([]byte, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func writeFile(filePath string, content string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
