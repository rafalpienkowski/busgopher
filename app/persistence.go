package main

import (
	"encoding/json"
	"io"
	"os"
)

func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
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

func loadConnections() ([]connection, error) {
	var connections []connection
	bytes, err := readFile("connections.json")
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
	bytes, err := readFile("messages.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil

}
