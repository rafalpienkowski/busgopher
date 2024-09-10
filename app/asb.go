package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func GetClient(namespace string) *azservicebus.Client {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azservicebus.NewClient(namespace, cred, nil)
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
		Body: []byte(message.body),
        Subject: &message.subject,
	}
	err = sender.SendMessage(context.TODO(), sbMessage, nil)
	if err != nil {
		panic(err)
	}
}
