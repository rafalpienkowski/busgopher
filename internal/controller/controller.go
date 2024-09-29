package controller

import (
	"github.com/rafalpienkowski/busgopher/internal/asb"
	"github.com/rafalpienkowski/busgopher/internal/config"
)

type Controller struct {
	Connections  []asb.Connection
	Messages     []asb.Message
	Destinations []string

	selectedConnection  asb.Connection
	selectedMessage     asb.Message
	selectedDestination string
}

func NewController(config *config.Config) (*Controller, error) {

	controller := Controller{}
	controller.Connections = *config.Connections
	controller.Messages = *config.Messages

	return &controller, nil
}

func (controller *Controller) SelectConnection(name string) {
	for _, conn := range controller.Connections {
		if conn.Name == name {
			controller.selectedConnection = conn
			controller.selectedDestination = ""
			return
		}
	}
}

func (controller *Controller) SelectMessage(name string) {
	for _, msg := range controller.Messages {
		if msg.Name == name {
			controller.selectedMessage = msg
			return
		}
	}
}

func (controller *Controller) SelectDestination(name string) {
	controller.selectedDestination = name
}

func (controller *Controller) Send() {
}
