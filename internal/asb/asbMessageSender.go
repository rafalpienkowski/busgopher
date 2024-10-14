package asb

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type AsbMessageSender struct {
	credentials *azidentity.DefaultAzureCredential
	client      *azservicebus.Client
}

func (sender *AsbMessageSender) getCredentials() error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}
	sender.credentials = cred

	return nil
}

func (sender *AsbMessageSender) getClient(namespace string) error {
	if credErr := sender.getCredentials(); credErr != nil {
		return credErr
	}

	client, err := azservicebus.NewClient(namespace, sender.credentials, nil)
	if err != nil {
		return err
	}

	sender.client = client

	return nil
}

func (messageSender *AsbMessageSender) Send(
	namespace string,
	destination string,
	message Message,
) error {
	if messageSender.client == nil {
		connErr := messageSender.getClient(namespace)
		if connErr != nil {
			return connErr
		}
	}
	sender, err := messageSender.client.NewSender(destination, nil)
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
