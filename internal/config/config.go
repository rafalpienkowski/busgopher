package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/rafalpienkowski/busgopher/internal/asb"
)

const connectionFilename = "connections.json"
const messagesFilename = "messages.json"

type Config struct {
	Connections *[]asb.Connection
	Messages    *[]asb.Message
}

func LoadConfig() *Config {

	config := Config{}

	connections, err := loadConnections()
	if err != nil {
		fmt.Println("Can't load connections")
		fmt.Println(err.Error())
		return nil
	}

	messages, err := loadMessages()
	if err != nil {
		fmt.Println("Can't load messages")
		fmt.Println(err.Error())
		return nil
	}

	config.Connections = &connections
	config.Messages = &messages

	return &config
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

func loadConnections() ([]asb.Connection, error) {
	var connections []asb.Connection

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

func loadMessages() ([]asb.Message, error) {
	var messages []asb.Message
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
