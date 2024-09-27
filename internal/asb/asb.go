package asb

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)


type busMessage struct {
    Name             string         `json:"name"`
	Body             string         `json:"body"`
	Subject          string         `json:"subject"`
	CustomProperties map[string]any `json:"customProperties"`
}

type busConnection struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
}

func GetClient(connection busConnection) *azservicebus.Client {
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
