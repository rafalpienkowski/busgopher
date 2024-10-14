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

func getTestConfig() *config.InMemoryConfigStorage {

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

func TestControllerShouldSetLoadedConfig(t *testing.T) {
	testConfig := getTestConfig()
	var givenConfig config.ConfigStorage = testConfig
	var buffer bytes.Buffer
	var writer io.Writer = &buffer

	sut, _ := NewController(givenConfig, writer)

	assert.Equal(t, testConfig.Config, sut.Config)
}

func TestControllerShouldSelectExistingConnectionByName(t *testing.T) {
	testConfig := getTestConfig()
	var givenConfig config.ConfigStorage = testConfig
	var buffer bytes.Buffer
	var writer io.Writer = &buffer
	sut, _ := NewController(givenConfig, writer)

	sut.SelectConnectionByName("test")

	assert.Equal(t, &(testConfig.Config.Connections)[0], sut.SelectedConnection)
}

func TestControllerShouldWriteErrorWhenSelectingNonExistingConnectionByName(t *testing.T) {
	testConfig := getTestConfig()
	var givenConfig config.ConfigStorage = testConfig
	var buffer bytes.Buffer
	var writer io.Writer = &buffer
	sut, _ := NewController(givenConfig, writer)

	sut.SelectConnectionByName("non-existing")

	assert.Equal(
		t,
		"[Error] Can't find connection with name: non-existing",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldSelectDestinationByName(t *testing.T) {
	testConfig := getTestConfig()
	var givenConfig config.ConfigStorage = testConfig
	var buffer bytes.Buffer
	var writer io.Writer = &buffer
	sut, _ := NewController(givenConfig, writer)
	sut.SelectConnectionByName("test")

	sut.SelectDestinationByName("queue")

	assert.Equal(t, "queue", sut.selectedDestination)
}

func TestControllerShouldWriteErrorWhenSelectingNontExistingQueueName(t *testing.T) {
	testConfig := getTestConfig()
	var givenConfig config.ConfigStorage = testConfig
	var buffer bytes.Buffer
	var writer io.Writer = &buffer
	sut, _ := NewController(givenConfig, writer)
	sut.SelectConnectionByName("test")

	sut.SelectDestinationByName("non-existing")

	assert.Equal(
		t,
		"[Error] Can't find destination with name: non-existing",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldWriteErrorWhenSelectingQueueWithoutSelectedConnection(t *testing.T) {
	testConfig := getTestConfig()
	var givenConfig config.ConfigStorage = testConfig
	var buffer bytes.Buffer
	var writer io.Writer = &buffer
	sut, _ := NewController(givenConfig, writer)

	sut.SelectDestinationByName("queue")

	assert.Equal(
		t,
		"[Error] Can't find destination with name: queue",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldSetMessageByName(t *testing.T) {
	testConfig := getTestConfig()
	var givenConfig config.ConfigStorage = testConfig
	var buffer bytes.Buffer
	var writer io.Writer = &buffer
	sut, _ := NewController(givenConfig, writer)

	sut.SelectMessageByName("test")

	assert.Equal(t, &(testConfig.Config.Messages)[0], sut.selectedMessage)
}

func TestControllerShouldWriteErrorWhenSelectingNonExistingMessage(t *testing.T) {
	testConfig := getTestConfig()
	var givenConfig config.ConfigStorage = testConfig
	var buffer bytes.Buffer
	var writer io.Writer = &buffer
	sut, _ := NewController(givenConfig, writer)

	sut.SelectMessageByName("non-existing")

	assert.Equal(
		t,
		"[Error] Can't find message with name: non-existing",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldNotSendWhenConnectionNotSelected(t *testing.T){
	testConfig := getTestConfig()
	var givenConfig config.ConfigStorage = testConfig
	var buffer bytes.Buffer
	var writer io.Writer = &buffer
	sut, _ := NewController(givenConfig, writer)

    sut.Send()

	assert.Equal(
		t,
		"[Error] Connection not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldNotSendWhenDestinationNotSelected(t *testing.T){
	testConfig := getTestConfig()
	var givenConfig config.ConfigStorage = testConfig
	var buffer bytes.Buffer
	var writer io.Writer = &buffer
	sut, _ := NewController(givenConfig, writer)
    sut.SelectConnectionByName("test")
    sut.SelectMessageByName("test")

    sut.Send()

	assert.Equal(
		t,
		"[Error] Destination not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func TestControllerShouldNotSendWhenMessageNotSelected(t *testing.T){
	testConfig := getTestConfig()
	var givenConfig config.ConfigStorage = testConfig
	var buffer bytes.Buffer
	var writer io.Writer = &buffer
	sut, _ := NewController(givenConfig, writer)
    sut.SelectConnectionByName("test")
    sut.SelectDestinationByName("queue")

    sut.Send()

	assert.Equal(
		t,
		"[Error] Message not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}
