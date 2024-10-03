package controller

import (
	"errors"
	"strings"

	"github.com/rafalpienkowski/busgopher/internal/asb"
	"github.com/rafalpienkowski/busgopher/internal/config"
)

type Controller struct {
	Connections []asb.Connection
	Messages    []asb.Message

	SelectedConnection  asb.Connection
	selectedMessage     asb.Message
	selectedDestination string
}

func NewController(config *config.Config) (*Controller, error) {

	controller := Controller{}
	controller.Connections = *config.Connections
	controller.Messages = *config.Messages

	return &controller, nil
}

func (controller *Controller) SelectConnectionByName(name string) error {
	for _, conn := range controller.Connections {
		if strings.EqualFold(conn.Name, name) {
			controller.SelectedConnection = conn
			return nil
		}
	}
	return errors.New("Can't find connection with name: " + name)
}

func (controller *Controller) SelectDestinationByName(name string) error {
	for _, dest := range controller.SelectedConnection.Entities {
		if strings.EqualFold(dest, name) {
			controller.selectedDestination = dest
			return nil
		}
	}
	return errors.New("Can't find destination with name: " + name)
}

func (controller *Controller) SelectMessageByName(name string) error {
	for _, msg := range controller.Messages {
		if strings.EqualFold(msg.Name, name) {
			controller.selectedMessage = msg
			return nil
		}
	}
	return errors.New("Can't find message with name: " + name)
}

func (controller *Controller) SelectMessage(message asb.Message) {
	controller.selectedMessage = message
}

func (controller *Controller) Send(destination string) error {
	return controller.SelectedConnection.SendMessage(destination, controller.selectedMessage)
}
