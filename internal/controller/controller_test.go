package controller

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafalpienkowski/busgopher/internal/asb"
	"github.com/rafalpienkowski/busgopher/internal/config"
)

func getInMemoryConfig() *config.InMemoryConfigStorage {

	return &config.InMemoryConfigStorage{
		Config: config.GetTestConfig(),
	}
}

func createTestController() (*Controller, *config.InMemoryConfigStorage, *asb.InMemoryMessageSender) {
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

	return controller, inMemoryConfig, inMemoryMessageSender
}

func Test_Controller_Should_Load_Config(t *testing.T) {
	controller, inMemoryConfig, _ := createTestController()

	assert.Equal(t, inMemoryConfig.Config, controller.Config)
}

func Test_Controller_Should_Select_Connection(t *testing.T) {
	controller, _, _ := createTestController()

	err := controller.SelectConnectionByName("test-connection")

	assert.NoError(t, err)
	assert.Equal(t,
		"test-connection",
		controller.selectedConnectionName)
}

func Test_Controller_Should_Write_Error_When_Selecting_NonExisting_Connection(t *testing.T) {
	controller, _, _ := createTestController()

	err := controller.SelectConnectionByName("non-existing")

	assert.Error(t, err, "Can't find connection with name: non-existing")
}

func Test_Controller_Should_Select_Destination(t *testing.T) {
	controller, _, _ := createTestController()
	err := controller.SelectConnectionByName("test-connection")
	assert.NoError(t, err)

	err = controller.SelectDestinationByName("queue")

	assert.NoError(t, err)
	assert.Equal(t, "queue", controller.selectedDestination)
}

func Test_Controller_Should_Write_Error_When_Selecting_Non_Existing_Destination(t *testing.T) {
	controller, _, _ := createTestController()
	err := controller.SelectConnectionByName("test-connection")
	assert.NoError(t, err)

	err = controller.SelectDestinationByName("non-existing")

	assert.Error(t, err, "Can't find destination with name: non-existing")
}

func Test_Controller_Should_Write_Error_When_Selecting_Queue_Without_Connection(t *testing.T) {
	controller, _, _ := createTestController()

	err := controller.SelectDestinationByName("queue")

	assert.Error(t, err, "Can't find destination with name: queue")
}

func Test_Controller_Should_Select_Message(t *testing.T) {
	controller, _, _ := createTestController()

	err := controller.SelectMessageByName("test-message")

	assert.NoError(t, err)
	assert.Equal(t, "test-message", controller.selectedMessageName)
}

func Test_Controller_Should_Write_Error_When_Selecting_NonExisting_Message(t *testing.T) {
	controller, _, _ := createTestController()

	err := controller.SelectMessageByName("non-existing")

	assert.Error(t, err, "Can't find message with name: non-existing")
}

func Test_Controller_Should_Not_Send_When_Connection_Not_Selected(t *testing.T) {
	controller, _, _ := createTestController()

	err := controller.Send()

	assert.Error(t, err, "Connection not selected!")
}

func Test_Controller_Should_Not_Send_When_Destination_Not_Selected(t *testing.T) {
	controller, _, _ := createTestController()
	err := controller.SelectConnectionByName("test-connection")
	assert.NoError(t, err)
	err = controller.SelectMessageByName("test-message")
	assert.NoError(t, err)

	err = controller.Send()

	assert.Error(t, err, "Destination not selected!")
}

func Test_Controller_Should_Not_Send_When_Message_Not_Selected(t *testing.T) {
	controller, _, _ := createTestController()
	err := controller.SelectConnectionByName("test-connection")
	assert.NoError(t, err)
	err = controller.SelectDestinationByName("queue")
	assert.NoError(t, err)

	err = controller.Send()

	assert.Error(t, err, "Message not selected!")
}

func Test_Controller_Should_Send_Message(t *testing.T) {
	controller, inMemoryConfig, messageSender := createTestController()
	err := controller.SelectConnectionByName("test-connection")
	assert.NoError(t, err)
	err = controller.SelectDestinationByName("queue")
	assert.NoError(t, err)
	err = controller.SelectMessageByName("test-message")
	assert.NoError(t, err)

	err = controller.Send()
	assert.NoError(t, err)

	assert.Equal(t, "test.azure.com", messageSender.Namespace)
	assert.Equal(t, "queue", messageSender.Destination)
	assert.Equal(t, inMemoryConfig.Config.Messages["test-message"], messageSender.Message)
}

func Test_Controller_Should_Get_Selected_Connection(t *testing.T) {
	controller, inMemoryConfig, _ := createTestController()
	err := controller.SelectConnectionByName("test-connection")
	assert.NoError(t, err)

	conn := controller.GetSelectedConnection()

	assert.Equal(t, inMemoryConfig.Config.Connections["test-connection"], *conn)
}

func Test_Controller_Should_Return_Nil_When_Connection_Not_Selected(t *testing.T) {
	controller, _, _ := createTestController()

	conn := controller.GetSelectedConnection()

	assert.Nil(t, conn)
}

func Test_Controller_Should_Return_Connections(t *testing.T) {
	controller, _, _ := createTestController()

	connections := controller.GetConnections()

	assert.Equal(t, []Connection{{
		Name:      "test-connection",
		Namespace: "test.azure.com",
	}}, connections)
}

func Test_Controller_Should_Return_Messages(t *testing.T) {
	controller, inMemoryConfig, _ := createTestController()

	messages := controller.GetMessages()

	assert.Equal(t, []Message{{
		Name:    "test-message",
		Message: inMemoryConfig.Config.Messages["test-message"],
	}}, messages)
}

func Test_Controller_Should_Return_No_Destinations_When_Connection_Not_Selected(t *testing.T) {
	controller, _, _ := createTestController()

	destinations := controller.GetDestiationNamesForSelectedConnection()

	assert.Equal(t, []string{}, destinations)
}

func Test_Controller_Should_Return_Destinations_For_Selected_Connection(t *testing.T) {
	controller, _, _ := createTestController()
	err := controller.SelectConnectionByName("test-connection")
	assert.NoError(t, err)

	destinations := controller.GetDestiationNamesForSelectedConnection()

	assert.Equal(t, []string{"queue", "topic"}, destinations)
}

func Test_Controller_Save_Config_Json_Should_Clear_Selections(t *testing.T) {
	controller, _, _ := createTestController()
	err := controller.SelectConnectionByName("test-connection")
	assert.NoError(t, err)
	err = controller.SelectDestinationByName("queue")
	assert.NoError(t, err)
	err = controller.SelectMessageByName("test-message")
	assert.NoError(t, err)

	err = controller.SaveConfigJson("{}")

	assert.NoError(t, err)
	assert.Equal(t, "", controller.selectedConnectionName)
	assert.Equal(t, "", controller.selectedDestination)
	assert.Equal(t, "", controller.selectedMessageName)
}

func Test_Controller_Save_Config_Json_Should_Save_Config(t *testing.T) {
	controller, inMemoryConfig, _ := createTestController()

	err := controller.SaveConfigJson("{}")

	assert.NoError(t, err)
	assert.Equal(t, config.Config{}, inMemoryConfig.Config)
}

func Test_Controller_Save_Config_Json_Should_Set_Controller_Config(t *testing.T) {
	controller, _, _ := createTestController()

	err := controller.SaveConfigJson("{}")

	assert.NoError(t, err)
	assert.Equal(t, config.Config{}, controller.Config)
}

func Test_Controller_Save_Config_Json_Should_Save_Something_More_Sophisticated(t *testing.T) {
	controller, _, _ := createTestController()

	err := controller.SaveConfigJson(`{
        "connections": {
            "another-connection": {
                "namespace": "another.azure.com",
                "destinations": [ "some-queue", "some-topic" ]
            }
        },
        "messages": {
            "another-message": {
                "messageID": "131",
                "correlationID": "22",
                "body": "{ another msg body }",
                "replyTo": "Please not",
                "subject": "Another",
                "customProperties": {
                    "first": "second",
                    "another": false
                }
            }
        }
    }`)

	assert.NoError(t, err)

	connections := make(map[string]asb.Connection)
	connections["another-connection"] = asb.Connection{
		Namespace: "another.azure.com",
		Destinations: []string{
			"some-queue",
			"some-topic",
		},
	}
	customProperties := make(map[string]any)
	customProperties["first"] = "second"
	customProperties["another"] = false

	messages := make(map[string]asb.Message)
	messages["another-message"] = asb.Message{
		Body:             "{ another msg body }",
		MessageID:        "131",
		CorrelationID:    "22",
		ReplayTo:         "Please not",
		Subject:          "Another",
		CustomProperties: customProperties,
	}

	expected := config.Config{
		Connections: connections,
		Messages:    messages,
	}
	assert.Equal(t, expected, controller.Config)
}
