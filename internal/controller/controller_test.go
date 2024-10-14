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
