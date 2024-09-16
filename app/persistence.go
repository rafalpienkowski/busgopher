package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const connectionFilename = "connections.json"
const messagesFilename = "messages.json"

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

func loadConnections() ([]connection, error) {
	var connections []connection

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
    
    fmt.Println("Bytes: ", len(bytes))
	err = json.Unmarshal(bytes, &messages)
	if err != nil {
		return nil, err
	}
    
	return messages, nil

}
