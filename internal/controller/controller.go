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
	Config        config.Config
	configStorage config.ConfigStorage

	SelectedConnection  *asb.Connection
	selectedMessage     *asb.Message
	selectedDestination string

	selectedConnectionName string
	selectedMessageName    string

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
	controller.configStorage = configStorage

	return &controller, nil
}

func (controller *Controller) SelectConnectionByName(name string) {

	conn, ok := controller.Config.NConnections[name]
	if ok {
		controller.selectedConnectionName = name
		controller.selectedDestination = ""
		controller.writeLog("Selected connection: " + conn.Name + " (" + conn.Namespace + ")")
		return
	}

	controller.writeError("Can't find connection with name: " + name)
}

func (controller *Controller) SelectDestinationByName(name string) {
	conn, ok := controller.Config.NConnections[controller.selectedConnectionName]
	if ok {
		for _, dest := range conn.Destinations {
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
	_, ok := controller.Config.NMessages[name]
	if ok {
		controller.selectedMessageName = name
		controller.writeLog("Selected message: " + name)
		return
	}
	controller.writeError("Can't find message with name: " + name)
}

func (controller *Controller) AddDestination(newDestination string) {
	if len(controller.selectedConnectionName) == 0 {
		controller.writeError("Connection not selected!")
		return
	}

	conn, ok := controller.Config.NConnections[controller.selectedConnectionName]
	if ok {

		conn.Destinations = append(
			conn.Destinations,
			newDestination,
		)

		controller.Config.NConnections[controller.selectedConnectionName] = conn
		controller.saveConfig()
		return
	}

	controller.writeError("Can't find selected connection" + controller.selectedConnectionName)
}

func (controller *Controller) RemoveDestination(destination string) {
	if controller.SelectedConnection == nil {
		controller.writeError("Connection not selected!")
		return
	}

	for i, conn := range controller.Config.Connections {
		if conn.Name == controller.SelectedConnection.Name {
			newDestinations := []string{}
			for _, d := range conn.Destinations {
				if d != destination {
					newDestinations = append(newDestinations, d)
				}
			}

			if len(newDestinations) == len(conn.Destinations) {
				controller.writeError("Nothing to remove!")
				return
			}

			controller.Config.Connections[i].Destinations = newDestinations

			controller.saveConfig()
			controller.SelectConnectionByName(controller.SelectedConnection.Name)
		}
	}
}

func (controller *Controller) UpdateDestination(oldDestination string, newDestination string) {
	if controller.SelectedConnection == nil {
		controller.writeError("Connection not selected!")
		return
	}

	for i, conn := range controller.Config.Connections {
		if conn.Name == controller.SelectedConnection.Name {
			newDestinations := []string{}
			for _, d := range conn.Destinations {
				if d == oldDestination {
					newDestinations = append(newDestinations, newDestination)
				} else {
					newDestinations = append(newDestinations, d)
				}
			}

			controller.Config.Connections[i].Destinations = newDestinations

			controller.saveConfig()
			controller.SelectConnectionByName(controller.SelectedConnection.Name)
		}
	}
}

func (controller *Controller) AddMessage(message asb.Message) {

	for _, msg := range controller.Config.Messages {
		if msg.Name == message.Name {
			controller.writeError("Message with name " + message.Name + " already exist")
			return
		}
	}

	controller.Config.Messages = append(controller.Config.Messages, message)

	controller.saveConfig()
}

func (controller *Controller) RemoveMessage(name string) {

	newMessages := []asb.Message{}
	for _, msg := range controller.Config.Messages {
		if msg.Name != name {
			newMessages = append(newMessages, msg)
		}
	}
	if len(newMessages) == len(controller.Config.Messages) {
		controller.writeError("No message to remove")
		return
	}

	controller.Config.Messages = newMessages

	controller.saveConfig()
}

func (controller *Controller) UpdateMessage(message asb.Message) {
	if controller.selectedMessage == nil {
		controller.writeError("Message not selected!")
		return
	}

	newMessages := []asb.Message{}
	for _, msg := range controller.Config.Messages {
		if msg.Name != controller.selectedMessage.Name {
			newMessages = append(newMessages, msg)
		} else {
			newMessages = append(newMessages, message)
		}
	}

	controller.Config.Messages = newMessages

	controller.saveConfig()
	controller.SelectMessageByName(message.Name)
}

func (controller *Controller) RemoveSelectedConnection() {
	if controller.SelectedConnection == nil {
		controller.writeError("Connection not selected!")
		return
	}

	newConnections := []asb.Connection{}
	for _, conn := range controller.Config.Connections {
		if conn.Name != controller.SelectedConnection.Name {
			newConnections = append(newConnections, conn)
		}
	}

	if len(newConnections) == len(controller.Config.Connections) {
		controller.writeError("No connection to remove")
		return
	}

	controller.Config.Connections = newConnections
	controller.saveConfig()
	controller.SelectedConnection = nil
}

func (controller *Controller) AddConnection(newConnection *asb.Connection) {
	newConnections := []asb.Connection{}
	for _, conn := range controller.Config.Connections {
		if conn.Name == newConnection.Name {
			controller.writeError("Connection '" + newConnection.Name + "' exist")
			return
		} else {
			newConnections = append(newConnections, conn)
		}
	}
	newConnections = append(newConnections, *newConnection)

	controller.Config.Connections = newConnections
	controller.saveConfig()
}

func (controller *Controller) UpdateSelectedConnection(newConnection asb.Connection) {
	if controller.SelectedConnection == nil {
		controller.writeError("Connection not selected!")
		return
	}

	newConnections := []asb.Connection{}
	for _, conn := range controller.Config.Connections {
		if conn.Name == controller.SelectedConnection.Name {
			newConnections = append(newConnections, newConnection)
		} else {
			newConnections = append(newConnections, conn)
		}
	}

	controller.Config.Connections = newConnections
	controller.saveConfig()
	controller.SelectConnectionByName(newConnection.Name)
}

func (controller *Controller) Send() {

	if len(controller.selectedConnectionName) == 0 {
		controller.writeError("Connection not selected!")
		return
	}

	if len(controller.selectedMessageName) == 0 {
		controller.writeError("Message not selected!")
		return
	}

	if len(controller.selectedDestination) == 0 {
		controller.writeError("Destination not selected!")
		return
	}

	controller.writeLog("Sending message to: " + controller.selectedDestination)
	controller.writeLog(
		"Sending message to: " + controller.Config.NConnections[controller.selectedConnectionName].Namespace,
	)

	err := controller.messageSender.Send(
		controller.Config.NConnections[controller.selectedConnectionName].Namespace,
		controller.selectedDestination,
		controller.Config.NMessages[controller.selectedMessageName],
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

func (controller *Controller) saveConfig() {
	err := controller.configStorage.Save(controller.Config)
	if err != nil {
		controller.writeError("Can't save config changes: " + err.Error())
	}
}
