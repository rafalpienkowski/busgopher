package controller

import (
	"github.com/rafalpienkowski/busgopher/internal/asb"
	"github.com/rafalpienkowski/busgopher/internal/config"
)

type Controller struct {
	Connections []asb.Connection
	Messages    []asb.Message
}

func NewController(config *config.Config) (*Controller, error) {

	controller := Controller{}
	controller.Connections = *config.Connections
	controller.Messages = *config.Messages

	return &controller, nil
}
