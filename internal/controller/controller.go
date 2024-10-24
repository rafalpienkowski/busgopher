package controller

import (
	"errors"
	"fmt"
	"strings"
	"time"

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

type Controller struct {
	Config        config.Config
	configStorage config.ConfigStorage

	selectedConnectionName string
	selectedMessageName    string
	selectedDestination    string

	messageSender asb.MessageSender
	print         Print
}

type Print func(string)

func NewController(
	configStorage config.ConfigStorage,
	messageSender asb.MessageSender,
	print Print,
) (*Controller, error) {

	config, err := configStorage.Load()
	if err != nil {
		return nil, err
	}

	controller := Controller{}
	controller.Config = config
	controller.messageSender = messageSender
	controller.configStorage = configStorage
	controller.print = print

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

	controller.writeLog("Sending message to: " + controller.selectedDestination)
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

func (controller *Controller) writeLog(log string) {
	controller.print(fmt.Sprintf(
		"[%v]: [Info] %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		log,
	))
}

func (controller *Controller) saveconfig() error {
	return controller.configStorage.Save(controller.Config)
}
