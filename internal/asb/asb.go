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
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
}

func GetClient(connection Connection) *azservicebus.Client {
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

func SendMessage(queueName string, message Message, client *azservicebus.Client) {
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
