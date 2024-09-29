package asb

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type Message struct {
	Name             string         `json:"name"`
	Body             string         `json:"body"`
	Subject          string         `json:"subject"`
	CustomProperties map[string]any `json:"customProperties"`
}

func (msg *Message) Print() string {

	prettyMsgBytes, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(prettyMsgBytes)
}

type Connection struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`

	credentials *azidentity.DefaultAzureCredential
	client      *azservicebus.Client
}

func (connection *Connection) getCredentials() error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}
	connection.credentials = cred

	return nil
}

func (connection *Connection) getClient() error {
	if credErr := connection.getCredentials(); credErr != nil {
		return credErr
	}

	client, err := azservicebus.NewClient(connection.Namespace, connection.credentials, nil)
	if err != nil {
		return err
	}

	connection.client = client

	return nil
}

func (connection *Connection) SendMessage(destination string, message Message) error {
	if connection.client == nil {
		connErr := connection.getClient()
		if connErr != nil {
			return connErr
		}
	}
	sender, err := connection.client.NewSender(destination, nil)
	if err != nil {
		return err
	}
	defer sender.Close(context.TODO())

	sbMessage := &azservicebus.Message{
		Body:                  []byte(message.Body),
		Subject:               &message.Subject,
		ApplicationProperties: message.CustomProperties,
	}
	err = sender.SendMessage(context.TODO(), sbMessage, nil)

	return err
}
