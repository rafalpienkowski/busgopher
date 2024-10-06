package asb

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type Connection struct {
	Name         string   `json:"name"`
	Namespace    string   `json:"namespace"`
	Destinations []string `json:"destinations"`

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

	body, err := message.TransformBody()
	if err != nil {
		return err
	}

	sbMessage := &azservicebus.Message{
		Body: []byte(body),
	}

	if message.CorrelationID != "" {
		sbMessage.CorrelationID = &message.CorrelationID
	}

	if message.MessageID != "" {
		sbMessage.MessageID = &message.MessageID
	}

	if message.ReplayTo != "" {
		sbMessage.ReplyTo = &message.ReplayTo
	}

	if message.Subject != "" {
		sbMessage.Subject = &message.Subject
	}

	if len(message.CustomProperties) > 0 {
		sbMessage.ApplicationProperties = message.CustomProperties
	}

	err = sender.SendMessage(context.TODO(), sbMessage, nil)

	return err
}
