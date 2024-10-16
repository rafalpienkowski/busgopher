package controller

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafalpienkowski/busgopher/internal/asb"
	"github.com/rafalpienkowski/busgopher/internal/config"
)

func getInMemoryConfig() *config.InMemoryConfigStorage {

	return &config.InMemoryConfigStorage{
		Config: config.Config{
			Connections: []asb.Connection{
				{
					Name:      "test",
					Namespace: "test.azure.com",
					Destinations: []string{
						"queue",
						"topic",
					},
				},
			},
			Messages: []asb.Message{
				{
					Name: "test",
					Body: "test msg body",
				},
			},
		},
	}
}

func getLastLine(log string) string {
	lines := strings.Split(log, "\n")
	if len(lines) > 1 {
		return lines[len(lines)-2]
	}
	return ""
}

func trimDatePart(log string) string {
	if len(log) < 23 {
		return ""
	}

	return log[23:]
}

func createTestController() (*Controller, *config.InMemoryConfigStorage, *asb.InMemoryMessageSender, *bytes.Buffer) {
	inMemoryConfig := getInMemoryConfig()
	var testConfig config.ConfigStorage = inMemoryConfig
	inMemoryMessageSender := &asb.InMemoryMessageSender{}
	var testMessageSender asb.MessageSender = inMemoryMessageSender
	var buffer bytes.Buffer
	var writer io.Writer = &buffer

	controller, _ := NewController(testConfig, testMessageSender, writer)

	return controller, inMemoryConfig, inMemoryMessageSender, &buffer
}

func TestControllerShouldSetLoadedConfig(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()

	assert.Equal(t, inMemoryConfig.Config, controller.Config)
}

func TestControllerShouldSelectExistingConnectionByName(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()

	controller.SelectConnectionByName("test")

	assert.Equal(t, &(inMemoryConfig.Config.Connections)[0], controller.SelectedConnection)
}

func TestControllerShouldWriteErrorWhenSelectingNonExistingConnectionByName(t *testing.T) {
	controller, _, _, buffer := createTestController()

	controller.SelectConnectionByName("non-existing")

	assert.Equal(
		t,
		"[Error] Can't find connection with name: non-existing",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldSelectDestinationByName(t *testing.T) {
	controller, _, _, _ := createTestController()
	controller.SelectConnectionByName("test")

	controller.SelectDestinationByName("queue")

	assert.Equal(t, "queue", controller.selectedDestination)
}

func TestControllerShouldWriteErrorWhenSelectingNontExistingQueueName(t *testing.T) {
	controller, _, _, buffer := createTestController()
	controller.SelectConnectionByName("test")

	controller.SelectDestinationByName("non-existing")

	assert.Equal(
		t,
		"[Error] Can't find destination with name: non-existing",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldWriteErrorWhenSelectingQueueWithoutSelectedConnection(t *testing.T) {
	controller, _, _, buffer := createTestController()

	controller.SelectDestinationByName("queue")

	assert.Equal(
		t,
		"[Error] Can't find destination with name: queue",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldSetMessageByName(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()

	controller.SelectMessageByName("test")

	assert.Equal(t, &(inMemoryConfig.Config.Messages)[0], controller.selectedMessage)
}

func TestControllerShouldWriteErrorWhenSelectingNonExistingMessage(t *testing.T) {
	controller, _, _, buffer := createTestController()

	controller.SelectMessageByName("non-existing")

	assert.Equal(
		t,
		"[Error] Can't find message with name: non-existing",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldNotSendWhenConnectionNotSelected(t *testing.T) {
	controller, _, _, buffer := createTestController()

	controller.Send()

	assert.Equal(
		t,
		"[Error] Connection not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldNotSendWhenDestinationNotSelected(t *testing.T) {
	controller, _, _, buffer := createTestController()
	controller.SelectConnectionByName("test")
	controller.SelectMessageByName("test")

	controller.Send()

	assert.Equal(
		t,
		"[Error] Destination not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldNotSendWhenMessageNotSelected(t *testing.T) {
	controller, _, _, buffer := createTestController()
	controller.SelectConnectionByName("test")
	controller.SelectDestinationByName("queue")

	controller.Send()

	assert.Equal(
		t,
		"[Error] Message not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldSendMessage(t *testing.T) {
	controller, inMemoryConfig, messageSender, _ := createTestController()
	controller.SelectConnectionByName("test")
	controller.SelectDestinationByName("queue")
	controller.SelectMessageByName("test")

	controller.Send()

	assert.Equal(t, "test.azure.com", messageSender.Namespace)
	assert.Equal(t, "queue", messageSender.Destination)
	assert.Equal(t, &(inMemoryConfig.Config.Messages)[0], &messageSender.Message)
}

func TestControllerShouldNotAddDestinationWhenConnectionNotSelected(t *testing.T) {
	controller, _, _, buffer := createTestController()

	controller.AddDestination("newDestination")

	assert.Equal(
		t,
		"[Error] Connection not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldAddDestinationWhenConnectionSelected(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	controller.SelectConnectionByName("test")

	controller.AddDestination("newDestination")

	assert.Equal(t, inMemoryConfig.Config,
		config.Config{
			Connections: []asb.Connection{
				{
					Name:      "test",
					Namespace: "test.azure.com",
					Destinations: []string{
						"queue",
						"topic",
						"newDestination",
					},
				},
			},
			Messages: []asb.Message{
				{
					Name: "test",
					Body: "test msg body",
				},
			},
		})
	assert.Equal(
		t,
		[]string{"queue", "topic", "newDestination"},
		controller.SelectedConnection.Destinations,
	)
}

func TestControllerShouldNotRemoveDestinationWhenConnectionNotSelected(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()

	controller.RemoveDestination("quque")

	assert.Equal(t, inMemoryConfig.Config,
		config.Config{
			Connections: []asb.Connection{
				{
					Name:      "test",
					Namespace: "test.azure.com",
					Destinations: []string{
						"queue",
						"topic",
					},
				},
			},
			Messages: []asb.Message{
				{
					Name: "test",
					Body: "test msg body",
				},
			},
		})
}

func TestControllerShouldNotRemoveDestinationWhenDestinationNotFound(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	controller.SelectConnectionByName("test")

	controller.RemoveDestination("notExisting")

	assert.Equal(t, inMemoryConfig.Config,
		config.Config{
			Connections: []asb.Connection{
				{
					Name:      "test",
					Namespace: "test.azure.com",
					Destinations: []string{
						"queue",
						"topic",
					},
				},
			},
			Messages: []asb.Message{
				{
					Name: "test",
					Body: "test msg body",
				},
			},
		})
	assert.Equal(
		t,
		[]string{"queue", "topic"},
		controller.SelectedConnection.Destinations,
	)
}

func TestControllerShouldRemoveDestinationWhenConnectionSelected(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	controller.SelectConnectionByName("test")

	controller.RemoveDestination("queue")

	assert.Equal(t,
		config.Config{
			Connections: []asb.Connection{
				{
					Name:      "test",
					Namespace: "test.azure.com",
					Destinations: []string{
						"topic",
					},
				},
			},
			Messages: []asb.Message{
				{
					Name: "test",
					Body: "test msg body",
				},
			},
		},
		inMemoryConfig.Config)

	assert.Equal(
		t,
		[]string{"topic"},
		controller.SelectedConnection.Destinations,
	)
}

func TestControllerShouldNotUpdateDestinationWhenConnectionNotSelected(t *testing.T) {
	controller, _, _, buffer := createTestController()

	controller.UpdateDestination("queue", "new-queue")

	assert.Equal(
		t,
		"[Error] Connection not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldNotUpdateDestinationWhenDestinationNotFound(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	controller.SelectConnectionByName("test")

	controller.UpdateDestination("non-existing-queue", "new-queue")

	assert.Equal(t,
		config.Config{
			Connections: []asb.Connection{
				{
					Name:      "test",
					Namespace: "test.azure.com",
					Destinations: []string{
						"queue",
						"topic",
					},
				},
			},
			Messages: []asb.Message{
				{
					Name: "test",
					Body: "test msg body",
				},
			},
		},
		inMemoryConfig.Config)

	assert.Equal(
		t,
		[]string{"queue", "topic"},
		controller.SelectedConnection.Destinations,
	)
}

func TestControllerShouldUpdateDestination(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	controller.SelectConnectionByName("test")

	controller.UpdateDestination("queue", "new-queue")

	assert.Equal(t,
		config.Config{
			Connections: []asb.Connection{
				{
					Name:      "test",
					Namespace: "test.azure.com",
					Destinations: []string{
						"new-queue",
						"topic",
					},
				},
			},
			Messages: []asb.Message{
				{
					Name: "test",
					Body: "test msg body",
				},
			},
		},
		inMemoryConfig.Config)

	assert.Equal(
		t,
		[]string{"new-queue", "topic"},
		controller.SelectedConnection.Destinations,
	)
}

func TestControllerShouldRemoveMessage(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()

	controller.RemoveMessage("test")

	assert.Equal(t,
		config.Config{
			Connections: []asb.Connection{
				{
					Name:      "test",
					Namespace: "test.azure.com",
					Destinations: []string{
						"queue",
						"topic",
					},
				},
			},
			Messages: []asb.Message{},
		},
		inMemoryConfig.Config)
}

func TestControllerShouldNotRemoveMessageWithUnknownName(t *testing.T) {
	controller, inMemoryConfig, _, buffer := createTestController()

	controller.RemoveMessage("unknown")

	assert.Equal(t,
		config.Config{
			Connections: []asb.Connection{
				{
					Name:      "test",
					Namespace: "test.azure.com",
					Destinations: []string{
						"queue",
						"topic",
					},
				},
			},
			Messages: []asb.Message{
				{
					Name: "test",
					Body: "test msg body",
				},
			},
		},
		inMemoryConfig.Config)

	assert.Equal(
		t,
		"[Error] No message to remove",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldAddNewMessage(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	newMsg := asb.Message{
		Name: "newMessage",
	}

	controller.AddMessage(newMsg)

	assert.Equal(t,
		config.Config{
			Connections: []asb.Connection{
				{
					Name:      "test",
					Namespace: "test.azure.com",
					Destinations: []string{
						"queue",
						"topic",
					},
				},
			},
			Messages: []asb.Message{
				{
					Name: "test",
					Body: "test msg body",
				},
				{
					Name: "newMessage",
				},
			},
		},
		inMemoryConfig.Config)
}

func TestControllerShouldNotAddNewMessageWhenNameIsNotUnique(t *testing.T) {
	controller, _, _, buffer := createTestController()
	newMsg := asb.Message{
		Name: "test",
	}

	controller.AddMessage(newMsg)

	assert.Equal(
		t,
		"[Error] Message with name test already exist",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldUpdateMessage(t *testing.T){
	controller, inMemoryConfig, _, _ := createTestController()
	newMsg := asb.Message{
		Name: "test",
        Body: "new msg body",
	}
    controller.SelectMessageByName("test")

	controller.UpdateMessage(newMsg)

	assert.Equal(t,
		config.Config{
			Connections: []asb.Connection{
				{
					Name:      "test",
					Namespace: "test.azure.com",
					Destinations: []string{
						"queue",
						"topic",
					},
				},
			},
			Messages: []asb.Message{
				{
					Name: "test",
					Body: "new msg body",
				},
			},
		},
		inMemoryConfig.Config)
}

func TestControllerShouldNotUpdateMessageWhenNoMessageSelected(t *testing.T){
	controller, _, _, buffer := createTestController()
	newMsg := asb.Message{
		Name: "test",
        Body: "new msg body",
	}

	controller.UpdateMessage(newMsg)

	assert.Equal(
		t,
		"[Error] Message not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}
