package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func GetClient(connection connection) *azservicebus.Client {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azservicebus.NewClient(connection.Namespace, cred, nil)
	if err != nil {
		panic(err)
	}
	return client
}

func SendMessage(queueName string, message busMessage, client *azservicebus.Client) {
	sender, err := client.NewSender(queueName, nil)
	if err != nil {
		panic(err)
	}
	defer sender.Close(context.TODO())

	sbMessage := &azservicebus.Message{
		Body: []byte(message.Body),
        Subject: &message.Subject,
        ApplicationProperties: message.CustomProperties,
	}
	err = sender.SendMessage(context.TODO(), sbMessage, nil)
	if err != nil {
		panic(err)
	}
}
