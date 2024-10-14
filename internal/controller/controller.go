package controller

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rafalpienkowski/busgopher/internal/asb"
	"github.com/rafalpienkowski/busgopher/internal/config"
)

type Controller struct {
	Config config.Config

	SelectedConnection  *asb.Connection
	selectedMessage     *asb.Message
	selectedDestination string

	messageSender asb.MessageSender
	logsWriter    io.Writer
}

func NewController(
	configStorage config.ConfigStorage,
	messageSender asb.MessageSender,
	logsWriter io.Writer,
) (*Controller, error) {

	config, err := configStorage.Load()
	if err != nil {
		return nil, err
	}

	controller := Controller{}
	controller.Config = config
	controller.messageSender = messageSender
	controller.logsWriter = logsWriter

	return &controller, nil
}

func (controller *Controller) SelectConnectionByName(name string) {
	for _, conn := range controller.Config.Connections {
		if strings.EqualFold(conn.Name, name) {
			controller.SelectedConnection = &conn
			controller.selectedDestination = ""
			controller.writeLog("Selected connection: " + conn.Name + " (" + conn.Namespace + ")")
			return
		}
	}
	controller.writeError("Can't find connection with name: " + name)
}

func (controller *Controller) SelectDestinationByName(name string) {
	if controller.SelectedConnection != nil {
		for _, dest := range controller.SelectedConnection.Destinations {
			if strings.EqualFold(dest, name) {
				controller.selectedDestination = dest
				controller.writeLog("Selected destination: " + name)
				return
			}
		}
	}
	controller.writeError("Can't find destination with name: " + name)
}

func (controller *Controller) SelectMessageByName(name string) {
	for _, msg := range controller.Config.Messages {
		if strings.EqualFold(msg.Name, name) {
			controller.selectedMessage = &msg
			controller.writeLog("Selected message: " + msg.Name)
			return
		}
	}
	controller.writeError("Can't find message with name: " + name)
}

func (controller *Controller) Send() {

	if controller.SelectedConnection == nil {
		controller.writeError("Connection not selected!")
		return
	}

	if controller.selectedMessage == nil {
		controller.writeError("Message not selected!")
		return
	}

	if len(controller.selectedDestination) == 0 {
		controller.writeError("Destination not selected!")
		return
	}

	controller.writeLog("Sending message to: " + controller.selectedDestination)
	err := controller.messageSender.Send(
		controller.SelectedConnection.Namespace,
		controller.selectedDestination,
		*controller.selectedMessage,
	)

	if err != nil {
		controller.writeError(err.Error())
	}
	controller.writeLog("Message send")
}

func (controller *Controller) writeLog(log string) {
	fmt.Fprintf(
		controller.logsWriter,
		"[%v]: [Info] %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		log,
	)
}

func (controller *Controller) writeError(log string) {
	fmt.Fprintf(
		controller.logsWriter,
		"[%v]: [Error] %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		log,
	)
}
