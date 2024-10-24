package controller

import (
	"bytes"
	"fmt"
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
		Namespace: "test.azure.com",
		Destinations: []string{
			"queue",
			"topic",
		},
	}
	nmessages := make(map[string]asb.Message)
	nmessages["test-message"] = asb.Message{
		Body: "test msg body",
	}

	return &config.InMemoryConfigStorage{
		Config: config.Config{
			Connections: nconnections,
			Messages:    nmessages,
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

	controller, _ := NewController(
		testConfig,
		testMessageSender,
		func(s string) { fmt.Fprintf(writer, "%v", s) },
	)

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

func Test_Controller_Should_Select_Destination(t *testing.T) {
	controller, _, _, _ := createTestController()
	controller.SelectConnectionByName("test-connection")

	controller.SelectDestinationByName("queue")

	assert.Equal(t, "queue", controller.selectedDestination)
}

func Test_Controller_Should_Write_Error_When_Selecting_Non_Existing_Destination(t *testing.T) {
	controller, _, _, buffer := createTestController()
	controller.SelectConnectionByName("test-connection")

	controller.SelectDestinationByName("non-existing")

	assert.Equal(
		t,
		"[Error] Can't find destination with name: non-existing",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func Test_Controller_Should_Write_Error_When_Selecting_Queue_Without_Connection(t *testing.T) {
	controller, _, _, buffer := createTestController()

	controller.SelectDestinationByName("queue")

	assert.Equal(
		t,
		"[Error] Can't find destination with name: queue",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func Test_Controller_Should_Select_Message(t *testing.T) {
	controller, _, _, _ := createTestController()

	controller.SelectMessageByName("test-message")

	assert.Equal(t, "test-message", controller.selectedMessageName)
}

func Test_Controller_Should_Write_Error_When_Selecting_NonExisting_Message(t *testing.T) {
	controller, _, _, buffer := createTestController()

	controller.SelectMessageByName("non-existing")

	assert.Equal(
		t,
		"[Error] Can't find message with name: non-existing",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func Test_Controller_Should_Not_Send_When_Connection_Not_Selected(t *testing.T) {
	controller, _, _, buffer := createTestController()

	controller.Send()

	assert.Equal(
		t,
		"[Error] Connection not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func Test_Controller_Should_Not_Send_When_Destination_Not_Selected(t *testing.T) {
	controller, _, _, buffer := createTestController()
	controller.SelectConnectionByName("test-connection")
	controller.SelectMessageByName("test-message")

	controller.Send()

	assert.Equal(
		t,
		"[Error] Destination not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func Test_Controller_Should_Not_Send_When_Message_Not_Selected(t *testing.T) {
	controller, _, _, buffer := createTestController()
	controller.SelectConnectionByName("test-connection")
	controller.SelectDestinationByName("queue")

	controller.Send()

	assert.Equal(
		t,
		"[Error] Message not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func Test_Controller_Should_Send_Message(t *testing.T) {
	controller, inMemoryConfig, messageSender, _ := createTestController()
	controller.SelectConnectionByName("test-connection")
	controller.SelectDestinationByName("queue")
	controller.SelectMessageByName("test-message")

	controller.Send()

	assert.Equal(t, "test.azure.com", messageSender.Namespace)
	assert.Equal(t, "queue", messageSender.Destination)
	assert.Equal(t, inMemoryConfig.Config.Messages["test-message"], messageSender.Message)
}
