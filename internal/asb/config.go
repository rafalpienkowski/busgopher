package asb

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const configName = "config.json"

type Config struct {
	Connections *[]Connection `json:"connections"`
	Messages    *[]Message    `json:"messages"`
}

func LoadConfig() *Config {

	config, err := loadConfig()
	if err != nil {
		fmt.Println("Can't load configuration")
		fmt.Println(err.Error())
		return nil
	}

	return config
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

func loadConfig() (*Config, error) {
	var config *Config

	bytes, err := readFile(configName)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		bytes = createDefaultConfig()
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func createDefaultConfig() []byte {
    return []byte(
    `{
        "Connections": [],
        "Messages": []
    }`)
}
