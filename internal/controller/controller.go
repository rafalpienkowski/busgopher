package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/rafalpienkowski/busgopher/internal/asb"
	"github.com/rafalpienkowski/busgopher/internal/config"
)

type Connection struct {
	Name      string
	Namespace string
}

type Message struct {
	Name    string
	Message asb.Message
}

type WriteLog func(log string)

type Controller struct {
	Config        config.Config
	configStorage config.ConfigStorage

	selectedConnectionName string
	selectedMessageName    string
	selectedDestination    string

	messageSender asb.MessageSender
	writeLog      WriteLog
}

func NewController(
	configStorage config.ConfigStorage,
	messageSender asb.MessageSender,
	writeLog WriteLog,
) (*Controller, error) {

	config, err := configStorage.Load()
	if err != nil {
		return nil, err
	}

	controller := Controller{}
	controller.Config = config
	controller.messageSender = messageSender
	controller.configStorage = configStorage
    controller.writeLog = writeLog

	return &controller, nil
}

func (controller *Controller) GetSelectedConnection() *asb.Connection {
	if len(controller.selectedConnectionName) == 0 {
		return nil
	}
	conn, ok := controller.Config.Connections[controller.selectedConnectionName]
	if ok {
		return &conn
	}

	return nil
}

func (controller *Controller) GetConnections() []Connection {
	var connections = []Connection{}
	for key := range controller.Config.Connections {
		connections = append(
			connections,
			Connection{Name: key, Namespace: controller.Config.Connections[key].Namespace},
		)
	}
	return connections
}

func (controller *Controller) GetMessages() []Message {
	var messages = []Message{}
	for key := range controller.Config.Messages {
		messages = append(messages, Message{Name: key, Message: controller.Config.Messages[key]})
	}

	return messages
}

func (controller *Controller) SelectConnectionByName(name string) error {

	_, ok := controller.Config.Connections[name]
	if ok {
		controller.selectedConnectionName = name
		controller.selectedDestination = ""
        controller.writeLog("Connection '" + name + "' selected")

		return nil
	}

	return errors.New("Can't find connection with name: " + name)
}

func (controller *Controller) SelectDestinationByName(name string) error {
	conn, ok := controller.Config.Connections[controller.selectedConnectionName]
	if ok {
		for _, dest := range conn.Destinations {
			if strings.EqualFold(dest, name) {
				controller.selectedDestination = dest
                controller.writeLog("Destination '" + name + "' selected")

				return nil
			}
		}
	}
	return errors.New("Can't find destination with name: " + name)
}

func (controller *Controller) SelectMessageByName(name string) error {
	_, ok := controller.Config.Messages[name]
	if ok {
		controller.selectedMessageName = name
        controller.writeLog("Message '" + name + "' selected")

		return nil
	}
	return errors.New("Can't find message with name: " + name)
}

func (controller *Controller) GetDestiationNamesForSelectedConnection() []string {
	if len(controller.selectedConnectionName) == 0 {
		return []string{}
	}

	return controller.Config.Connections[controller.selectedConnectionName].Destinations
}

func (controller *Controller) Send() error {

	if len(controller.selectedConnectionName) == 0 {
		return errors.New("Connection not selected!")
	}

	if len(controller.selectedMessageName) == 0 {
		return errors.New("Message not selected!")
	}

	if len(controller.selectedDestination) == 0 {
		return errors.New("Destination not selected!")
	}

	controller.writeLog(
		"Sending message to: " + controller.Config.Connections[controller.selectedConnectionName].Namespace,
	)

	err := controller.messageSender.Send(
		controller.Config.Connections[controller.selectedConnectionName].Namespace,
		controller.selectedDestination,
		controller.Config.Messages[controller.selectedMessageName],
	)

	if err != nil {
		return err
	}
	controller.writeLog("Message send")
	return nil
}

func (controller *Controller) SaveConfigJson(configJson string) error {
	config := config.Config{}
    err := json.Unmarshal([]byte(configJson), &config)
    if (err != nil){
        return err
    }
	controller.Config = config

	controller.selectedConnectionName = ""
	controller.selectedDestination = ""
	controller.selectedMessageName = ""
    controller.writeLog("Config saved")

	return controller.configStorage.Save(controller.Config)
}

func (controller *Controller) GetConfigString() (string, error) {
	decoded, err := json.Marshal(controller.Config)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	errIndent := json.Indent(&out, decoded, "", "\t")
	return out.String(), errIndent
}
