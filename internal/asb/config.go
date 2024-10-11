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
        tst := err.Error()
		fmt.Println("Can't load configuration")
		fmt.Println(tst)
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

func writeFile(filePath string, content string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
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

func loadConfig() (*Config, error) {
	var config *Config

	bytes, err := readFile(configName)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		defaultFileContent := createDefaultConfigString()
		err = writeFile(configName, defaultFileContent)
		if err != nil {
			return nil, err
		}

		bytes = []byte(defaultFileContent)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func createDefaultConfigString() string {
	return `{
        "Connections": [],
        "Messages": []
    }`
}
