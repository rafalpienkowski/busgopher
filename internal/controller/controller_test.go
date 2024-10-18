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
			NConnections: nconnections,
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
	assert.Equal(t, inMemoryConfig.Config.NMessages["test-message"], messageSender.Message)
}

func Test_Controller_Should_Not_Add_Destination_When_Connection_Not_Selected(t *testing.T) {
	controller, _, _, buffer := createTestController()

	controller.AddDestination("newDestination")

	assert.Equal(
		t,
		"[Error] Connection not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func Test_Controller_Should_Add_Destination_When_Connection_Selected(t *testing.T) {
	controller, _, _, _ := createTestController()
	controller.SelectConnectionByName("test-connection")

	controller.AddDestination("newDestination")

	assert.Equal(
		t,
		[]string{"queue", "topic", "newDestination"},
		controller.Config.NConnections["test-connection"].Destinations,
	)
}

func Test_Controller_Should_Save_Config_After_Adding_Destination(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	controller.SelectConnectionByName("test-connection")

	controller.AddDestination("newDestination")

	nconnections := make(map[string]asb.Connection)
	nconnections["test-connection"] = asb.Connection{
		Name:      "test-connection",
		Namespace: "test.azure.com",
		Destinations: []string{
			"queue",
			"topic",
			"newDestination",
		},
	}
	assert.Equal(t, nconnections, inMemoryConfig.Config.NConnections)
}

func Test_Controller_Should_Not_Remove_Destination_When_Connection_Not_Selected(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()

	controller.RemoveDestination("queue")

	nconnections := make(map[string]asb.Connection)
	nconnections["test-connection"] = asb.Connection{
		Name:      "test-connection",
		Namespace: "test.azure.com",
		Destinations: []string{
			"queue",
			"topic",
		},
	}

	assert.Equal(t, nconnections, inMemoryConfig.Config.NConnections)
}

func Test_Controller_Should_Return_Error_When_Removing_Destination_Without_Selected_Connection(
	t *testing.T,
) {

	controller, _, _, buffer := createTestController()

	controller.RemoveDestination("queue")

	assert.Equal(
		t,
		"[Error] Connection not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func Test_Controller_Should_Not_Remove_Destination_When_Destination_NotFound(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	controller.SelectConnectionByName("test-connection")

	controller.RemoveDestination("notExisting")

	nconnections := make(map[string]asb.Connection)
	nconnections["test-connection"] = asb.Connection{
		Name:      "test-connection",
		Namespace: "test.azure.com",
		Destinations: []string{
			"queue",
			"topic",
		},
	}
	assert.Equal(t, nconnections, inMemoryConfig.Config.NConnections)
}

func Test_Controller_Should_Write_Error_When_Remove_Non_Existing_Destination(t *testing.T) {
	controller, _, _, buffer := createTestController()
	controller.SelectConnectionByName("test-connection")

	controller.RemoveDestination("notExisting")

	assert.Equal(
		t,
		"[Error] Nothing to remove!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func Test_Controller_Should_Remove_Destination(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	controller.SelectConnectionByName("test-connection")

	controller.RemoveDestination("queue")

	nconnections := make(map[string]asb.Connection)
	nconnections["test-connection"] = asb.Connection{
		Name:      "test-connection",
		Namespace: "test.azure.com",
		Destinations: []string{
			"queue",
			"topic",
		},
	}
	assert.Equal(t, nconnections, inMemoryConfig.Config.NConnections)
}

func Test_Controller_Should_Not_Update_Destination_When_Connection_Not_Selected(t *testing.T) {
	controller, _, _, buffer := createTestController()

	controller.UpdateDestination("queue", "new-queue")

	assert.Equal(
		t,
		"[Error] Connection not selected!",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func Test_Controller_Should_Not_Update_Destination_When_Destination_Not_Found(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	controller.SelectConnectionByName("test-connection")

	controller.UpdateDestination("non-existing-queue", "new-queue")

	nconnections := make(map[string]asb.Connection)
	nconnections["test-connection"] = asb.Connection{
		Name:      "test-connection",
		Namespace: "test.azure.com",
		Destinations: []string{
			"queue",
			"topic",
		},
	}
	assert.Equal(t, nconnections, inMemoryConfig.Config.NConnections)
}

func Test_Controller_Should_Update_Destination(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	controller.SelectConnectionByName("test-connection")

	controller.UpdateDestination("queue", "new-queue")

	nconnections := make(map[string]asb.Connection)
	nconnections["test-connection"] = asb.Connection{
		Name:      "test-connection",
		Namespace: "test.azure.com",
		Destinations: []string{
			"new-queue",
			"topic",
		},
	}
	assert.Equal(t, nconnections, inMemoryConfig.Config.NConnections)
}

func Test_Controller_Should_Remove_Message(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()

	controller.RemoveMessage("test-message")

	nmessages := make(map[string]asb.Message)

	assert.Equal(t, nmessages, inMemoryConfig.Config.NMessages)
}

func TestControllerShouldNotRemoveMessageWithUnknownName(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()

	controller.RemoveMessage("unknown")

	nmessages := make(map[string]asb.Message)
	nmessages["test-message"] = asb.Message{
		Name: "test-message",
		Body: "test msg body",
	}
	assert.Equal(t, nmessages, inMemoryConfig.Config.NMessages)
}

func TestControllerShouldAddNewMessage(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	newMsg := asb.Message{
		Name: "new-message",
		Body: "new msg body",
	}

	controller.AddMessage(newMsg)

	nmessages := make(map[string]asb.Message)
	nmessages["test-message"] = asb.Message{
		Name: "test-message",
		Body: "test msg body",
	}
	nmessages["new-message"] = newMsg
	assert.Equal(t, nmessages, inMemoryConfig.Config.NMessages)
}

func Test_Controller_Should_Not_Add_New_Message_When_Name_Is_Not_Unique(t *testing.T) {
	controller, _, _, buffer := createTestController()
	newMsg := asb.Message{
		Name: "test-message",
	}

	controller.AddMessage(newMsg)

	assert.Equal(
		t,
		"[Error] Message with name test-message already exist",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func Test_Controller_Should_Update_Message(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	newMsg := asb.Message{
		Name: "test-message",
		Body: "new msg body",
	}

	controller.UpdateMessage(newMsg)

	nmessages := make(map[string]asb.Message)
	nmessages["test-message"] = newMsg
	assert.Equal(t, nmessages, inMemoryConfig.Config.NMessages)
}

func Test_Controller_Should_Not_Update_Message_When_Message_Not_Found(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	newMsg := asb.Message{
		Name: "test-message",
		Body: "test msg body",
	}

	controller.UpdateMessage(newMsg)

	nmessages := make(map[string]asb.Message)
	nmessages["test-message"] = newMsg
	assert.Equal(t, nmessages, inMemoryConfig.Config.NMessages)
}

func Test_Controller_Should_Write_Error_When_Update_Non_Existing_Message(t *testing.T) {
	controller, _, _, buffer := createTestController()
	newMsg := asb.Message{
		Name: "non-existing",
		Body: "test msg body",
	}

	controller.UpdateMessage(newMsg)

	assert.Equal(
		t,
		"[Error] Message with name non-existing not found",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func Test_Controller_Should_Remove_Connection(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()

	controller.RemoveConnection("test-connection")

	nconnections := make(map[string]asb.Connection)
	assert.Equal(t, nconnections, inMemoryConfig.Config.NConnections)
}

func Test_Controller_Should_Not_Remove_Non_Existing_Connection(t *testing.T) {
	controller, _, _, buffer := createTestController()

	controller.RemoveConnection("non-existing-connection")

	assert.Equal(
		t,
		"[Error] Connection not found",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func Test_Controller_Should_Clear_Selected_Connection_Name_When_Connection_Was_Removed(
	t *testing.T,
) {
	controller, _, _, _ := createTestController()
	controller.SelectConnectionByName("test-connection")
	assert.Equal(t, "test-connection", controller.selectedConnectionName)

	controller.RemoveConnection("test-connection")

	assert.Equal(t, "", controller.selectedConnectionName)
}

func Test_Controller_Should_Add_Connection(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	newConn := asb.Connection{
		Name:      "new-connection",
		Namespace: "new.azure.com",
	}

	controller.AddConnection(&newConn)

	nconnections := make(map[string]asb.Connection)
	nconnections["test-connection"] = asb.Connection{
		Name:      "test-connection",
		Namespace: "test.azure.com",
		Destinations: []string{
			"queue",
			"topic",
		},
	}
	nconnections["new-connection"] = asb.Connection{
		Name:      "new-connection",
		Namespace: "new.azure.com",
	}

	assert.Equal(t, nconnections, inMemoryConfig.Config.NConnections)
}

func Test_Controller_Should_Not_Add_New_Connection_When_Name_Exist(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	newConn := asb.Connection{
		Name:      "test-connection",
		Namespace: "test.azure.com",
	}

	controller.AddConnection(&newConn)

	nconnections := make(map[string]asb.Connection)
	nconnections["test-connection"] = asb.Connection{
		Name:      "test-connection",
		Namespace: "test.azure.com",
		Destinations: []string{
			"queue",
			"topic",
		},
	}
	assert.Equal(t, nconnections, inMemoryConfig.Config.NConnections)
}

func Test_Controller_Should_Write_Error_When_Add_New_Connection_And_Name_Exist(t *testing.T) {
	controller, _, _, buffer := createTestController()
	newConn := asb.Connection{
		Name:      "test-connection",
		Namespace: "new.azure.com",
	}

	controller.AddConnection(&newConn)

	assert.Equal(
		t,
		"[Error] Connection 'test-connection' exist",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}

func Test_Controller_Should_Update_Connection(t *testing.T) {
	controller, inMemoryConfig, _, _ := createTestController()
	newConn := asb.Connection{
		Name:      "new-connection",
		Namespace: "newtest.azure.com",
		Destinations: []string{
			"new-queue",
		},
	}
	controller.SelectConnectionByName("test-connection")
	controller.UpdateSelectedConnection(newConn)

	nconnections := make(map[string]asb.Connection)
	nconnections["new-connection"] = asb.Connection{
		Name:      "new-connection",
		Namespace: "newtest.azure.com",
		Destinations: []string{
			"new-queue",
		},
	}
	assert.Equal(t, nconnections, inMemoryConfig.Config.NConnections)
}

func Test_Controller_Should_Select_New_Connection(t *testing.T) {
	controller, _, _, _ := createTestController()
	newConn := asb.Connection{
		Name:      "new-connection",
		Namespace: "newtest.azure.com",
		Destinations: []string{
			"new-queue",
		},
	}
	controller.SelectConnectionByName("test-connection")
	controller.UpdateSelectedConnection(newConn)

	assert.Equal(t, "new-connection", controller.selectedConnectionName)
}

func Test_Controller_Should_Write_Error_When_Connection_For_Update_Not_Selected(t *testing.T){
	controller, _, _, buffer := createTestController()
	newConn := asb.Connection{
		Name:      "new-connection",
		Namespace: "newtest.azure.com",
		Destinations: []string{
			"new-queue",
		},
	}

	controller.UpdateSelectedConnection(newConn)

	assert.Equal(
		t,
		"[Error] Connection not selected",
		(trimDatePart(getLastLine(buffer.String()))),
	)
}
