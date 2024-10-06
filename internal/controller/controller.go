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
	Connections []asb.Connection
	Messages    []asb.Message

	SelectedConnection  *asb.Connection
	selectedMessage     *asb.Message
	selectedDestination string

	logsWriter io.Writer
}

func NewController(config *config.Config) (*Controller, error) {

	controller := Controller{}
	controller.Connections = *config.Connections
	controller.Messages = *config.Messages

	return &controller, nil
}

func (controller *Controller) SetLogsWriter(writer io.Writer) {
	controller.logsWriter = writer
}

func (controller *Controller) SelectConnectionByName(name string) {
	for _, conn := range controller.Connections {
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
	for _, dest := range controller.SelectedConnection.Destinations {
		if strings.EqualFold(dest, name) {
			controller.selectedDestination = dest
			controller.writeLog("Selected destination: " + name)
            return
		}
	}
	controller.writeError("Can't find destination with name: " + name)
}

func (controller *Controller) SelectMessageByName(name string) {
	for _, msg := range controller.Messages {
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
	err := controller.SelectedConnection.SendMessage(
		controller.selectedDestination,
		*controller.selectedMessage,
	)

	if err != nil {
		controller.writeError(err.Error())
	}
	controller.writeLog("Message send")
}

func (controller *Controller) writeLog(log string) {
	fmt.Fprintf(controller.logsWriter, "[%v]: %v\n", time.Now().Format("2006-01-02 15:04:05"), log)
}

func (controller *Controller) writeError(log string) {
	fmt.Fprintf(
		controller.logsWriter,
		"[%v]: [Error] %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		log,
	)
}
