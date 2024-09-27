package config

import (
	"io"
	"os"
)

const connectionFilename = "connections.json"
const messagesFilename = "messages.json"

type Config struct {

}

func LoadConfig() *Config {
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
/*
func loadConnections() ([]busConnection, error) {
	var connections []busConnection

	bytes, err := readFile(connectionFilename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &connections)
	if err != nil {
		return nil, err
	}

	return connections, nil
}

func loadMessages() ([]busMessage, error) {
	var messages []busMessage
	bytes, err := readFile(messagesFilename)
	if err != nil {
		return nil, err
	}

    if len(bytes) == 0 {
        bytes = []byte("[]")
    }
    
	err = json.Unmarshal(bytes, &messages)
	if err != nil {
		return nil, err
	}
    
	return messages, nil

}
*/
