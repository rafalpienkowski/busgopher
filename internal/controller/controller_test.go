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
	nconnections := make(map[string]asb.Connection)
	nconnections["test-connection"] = asb.Connection{
		Name:      "test-connection",
		Namespace: "test.azure.com",
		Destinations: []string{
			"queue",
			"topic",
		},
	}
	nmessages := make(map[string]asb.Message)
	nmessages["test-message"] = asb.Message{
		Name: "test-message",
		Body: "test msg body",
	}

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
			NConnections: nconnections,
			Messages: []asb.Message{
				{
					Name: "test",
					Body: "test msg body",
				},
			},
			NMessages: nmessages,
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

func Test_Controller_Should_Load_Config(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()

	assert.Equal(t, inMemoryConfig.Config, controller.Config)
}

func Test_Controller_Should_Select_Connection(t *testing.T) {
	controller, _, _, _ := createTestController()

	controller.SelectConnectionByName("test-connection")

	assert.Equal(t,
		"test-connection",
		controller.selectedConnectionName)
}

func Test_Controller_Should_Write_Error_When_Selecting_NonExisting_Connection(t *testing.T) {
	controller, _, _, buffer := createTestController()

	controller.SelectConnectionByName("non-existing")

	assert.Equal(
		t,
		"[Error] Can't find connection with name: non-existing",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

/*
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
*/

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

/*
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

func TestControllerShouldUpdateMessage(t *testing.T) {
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

func TestControllerShouldNotUpdateMessageWhenNoMessageSelected(t *testing.T) {
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

func TestControllerShouldRemoveSelectedConnection(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	controller.SelectConnectionByName("test")

	controller.RemoveSelectedConnection()

	assert.Equal(t,
		config.Config{
			Connections: []asb.Connection{},
			Messages: []asb.Message{
				{
					Name: "test",
					Body: "test msg body",
				},
			},
		},
		inMemoryConfig.Config)
}

func TestControllerShouldNotRemoveWhenConnectionNotSelected(t *testing.T) {
	controller, _, _, buffer := createTestController()

	controller.RemoveSelectedConnection()

	assert.Equal(
		t,
		"[Error] Connection not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldAddNewConnection(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	newConn := asb.Connection{
		Name:      "new connection",
		Namespace: "newtest.azure.com",
	}

	controller.AddConnection(&newConn)

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
				{
					Name:      "new connection",
					Namespace: "newtest.azure.com",
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
}

func TestControllerShouldNotAddNewConnectionWhenNameExist(t *testing.T) {
	controller, _, _, buffer := createTestController()
	newConn := asb.Connection{
		Name:      "test",
		Namespace: "newtest.azure.com",
	}

	controller.AddConnection(&newConn)

	assert.Equal(
		t,
		"[Error] Connection 'test' exist",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldNotUpdateConnectionWhenNoConnectionSelected(t *testing.T) {
	controller, _, _, buffer := createTestController()
	newConn := asb.Connection{
		Name:      "test",
		Namespace: "newtest.azure.com",
	}

	controller.UpdateSelectedConnection(newConn)

	assert.Equal(
		t,
		"[Error] Connection not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldUpdateSelectedConnection(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	newConn := asb.Connection{
		Name:      "new connection",
		Namespace: "newtest.azure.com",
	}
	controller.SelectConnectionByName("test")

	controller.UpdateSelectedConnection(newConn)

	assert.Equal(t,
		config.Config{
			Connections: []asb.Connection{
				{
					Name:      "new connection",
					Namespace: "newtest.azure.com",
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
}
*/
