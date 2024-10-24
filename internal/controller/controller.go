package controller

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rafalpienkowski/busgopher/internal/asb"
	"github.com/rafalpienkowski/busgopher/internal/config"
)

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

func (controller *Controller) SelectConnectionByName(name string) {

	conn, ok := controller.Config.Connections[name]
	if ok {
		controller.selectedConnectionName = name
		controller.selectedDestination = ""
		controller.writeLog("Selected connection: " + conn.Name + " (" + conn.Namespace + ")")
		return
	}

	controller.writeError("Can't find connection with name: " + name)
}

func (controller *Controller) SelectDestinationByName(name string) {
	conn, ok := controller.Config.Connections[controller.selectedConnectionName]
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
	_, ok := controller.Config.Messages[name]
	if ok {
		controller.selectedMessageName = name
		controller.writeLog("Selected message: " + name)
		return
	}
	controller.writeError("Can't find message with name: " + name)
}

func (controller *Controller) GetDestiationNamesForSelectedConnection() []string {
	if len(controller.selectedConnectionName) == 0 {
		return []string{}
	}

	return controller.Config.Connections[controller.selectedConnectionName].Destinations
}

func (controller *Controller) AddDestination(newDestination string) {
	if len(controller.selectedConnectionName) == 0 {
		controller.writeError("Connection not selected!")
		return
	}

	conn, ok := controller.Config.Connections[controller.selectedConnectionName]
	if ok {

		conn.Destinations = append(
			conn.Destinations,
			newDestination,
		)

		controller.Config.Connections[controller.selectedConnectionName] = conn
		controller.saveconfig()
		return
	}

	controller.writeError("Can't find selected connection" + controller.selectedConnectionName)
}

func (controller *Controller) RemoveDestination(destination string) {
	if len(controller.selectedConnectionName) == 0 {
		controller.writeError("Connection not selected!")
		return
	}

	conn, ok := controller.Config.Connections[controller.selectedConnectionName]
	if !ok {
		controller.writeError("Connection not found!")
		return
	}

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

	controller.Config.Connections[controller.selectedConnectionName] = conn

	controller.saveconfig()
}

func (controller *Controller) UpdateDestination(oldDestination string, newDestination string) {
	if len(controller.selectedConnectionName) == 0 {
		controller.writeError("Connection not selected!")
		return
	}

	conn, ok := controller.Config.Connections[controller.selectedConnectionName]
	if !ok {
		controller.writeError("Connection not found!")
		return
	}

	newDestinations := []string{}
	for _, d := range conn.Destinations {
		if d == oldDestination {
			newDestinations = append(newDestinations, newDestination)
		} else {
			newDestinations = append(newDestinations, d)
		}
	}
	conn.Destinations = newDestinations

	controller.Config.Connections[controller.selectedConnectionName] = conn

	controller.saveconfig()
}

func (controller *Controller) AddMessage(message asb.Message) error {
	_, ok := controller.Config.Messages[message.Name]
	if ok {
		controller.writeError("Message with name " + message.Name + " already exist")
		return errors.New("Message with name " + message.Name + " already exist") 
	}
	controller.Config.Messages[message.Name] = message

	controller.saveconfig()
    return nil
}

func (controller *Controller) RemoveMessage(name string) {
	delete(controller.Config.Messages, name)

	controller.selectedMessageName = ""

	controller.saveconfig()
}

func (controller *Controller) UpdateMessage(message asb.Message) {
	_, ok := controller.Config.Messages[message.Name]
	if !ok {
		controller.writeError("Message with name " + message.Name + " not found")
		return
	}
	controller.Config.Messages[message.Name] = message

	controller.saveconfig()
}

func (controller *Controller) RemoveSelectedConnection() (string, error) {

	name := controller.selectedConnectionName

	if len(name) == 0 {
		return "", errors.New("Select connection first!")
	}

	_, ok := controller.Config.Connections[name]
	if !ok {
		return "", errors.New("Connection not found")
	}

	delete(controller.Config.Connections, name)

	controller.saveconfig()
	controller.selectedConnectionName = ""

	return name, nil
}

func (controller *Controller) AddConnection(newConnection *asb.Connection) error {

	_, ok := controller.Config.Connections[newConnection.Name]
	if ok {
		return errors.New("Connection '" + newConnection.Name + "' exist")
	}

	controller.Config.Connections[newConnection.Name] = *newConnection
	controller.saveconfig()
	return nil
}

func (controller *Controller) UpdateSelectedConnection(newConnection *asb.Connection) error {

	if len(controller.selectedConnectionName) == 0 {
		controller.writeError("Connection not selected")
		return errors.New("Connection not selected!")
	}

	_, ok := controller.Config.Connections[controller.selectedConnectionName]
	if !ok {
		controller.writeError("Connection '" + controller.selectedConnectionName + "' not exist")
		return errors.New("Connection '" + controller.selectedConnectionName + "' not exist")
	}
	delete(controller.Config.Connections, controller.selectedConnectionName)

	controller.Config.Connections[newConnection.Name] = *newConnection
	controller.saveconfig()
	controller.SelectConnectionByName(newConnection.Name)
    return nil
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
		"Sending message to: " + controller.Config.Connections[controller.selectedConnectionName].Namespace,
	)

	err := controller.messageSender.Send(
		controller.Config.Connections[controller.selectedConnectionName].Namespace,
		controller.selectedDestination,
		controller.Config.Messages[controller.selectedMessageName],
	)

	if err != nil {
		controller.writeError(err.Error())
	}
	controller.writeLog("Message send")
}

func (controller *Controller) writeLog(log string) {
	controller.print(fmt.Sprintf(
		"[%v]: [Info] %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		log,
	))
}

func (controller *Controller) writeError(log string) {
	controller.print(fmt.Sprintf(
		"[%v]: [Error] %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		log,
	))
}

func (controller *Controller) saveconfig() {
	err := controller.configStorage.Save(controller.Config)
	if err != nil {
		controller.writeError("Can't save config changes: " + err.Error())
	}
}
