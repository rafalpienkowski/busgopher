package controller

import (
	"github.com/rafalpienkowski/busgopher/internal/asb"
	"github.com/rafalpienkowski/busgopher/internal/config"
)

type Controller struct {
	Connections []asb.Connection
	Messages    []asb.Message

	SelectedConnection asb.Connection
	selectedMessage    asb.Message
}

func NewController(config *config.Config) (*Controller, error) {

	controller := Controller{}
	controller.Connections = *config.Connections
	controller.Messages = *config.Messages

	return &controller, nil
}

func (controller *Controller) SelectMessage(message asb.Message) {
	controller.selectedMessage = message
}

func (controller *Controller) Send(destination string) error {
	return controller.SelectedConnection.SendMessage(destination, controller.selectedMessage)
}
