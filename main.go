package main

import (
	"fmt"
	"os"

    "github.com/rafalpienkowski/busgopher/internal/ui"
    "github.com/rafalpienkowski/busgopher/internal/config"
    "github.com/rafalpienkowski/busgopher/internal/controller"
)

func main() {

    config := config.LoadConfig()
    controller, err := controller.NewController(config)
	if err != nil {
		fmt.Printf("Failed to start: %v\n", err)
		os.Exit(1)
	}

    ui := ui.NewUI(controller)
    err = ui.Start()
	if err != nil {
		fmt.Printf("failed to start: %v\n", err)
		os.Exit(1)
	}

    /*
	sendButton := tview.NewButton("Send")
	sendButton.SetBorder(true)
	closeButton := tview.NewButton("Exit").SetSelectedFunc(func() {
		app.Stop()
	})
	closeButton.SetBorder(true)

	inputs := []tview.Primitive{
		connectionsSelection,
		messagesSelection,
		sendButton,
		closeButton,
	}

	connections, err := loadConnections()
	if err != nil {
		fmt.Fprintln(logs, "Can't load connections")
		fmt.Fprintln(logs, err.Error())
		return
	}

	for _, conn := range connections {
		connectionsSelection.AddItem(conn.Name, conn.Namespace, 'a', func() {
			fmt.Fprintln(logs, "Selected "+conn.Name)
			activeConnection = conn
		})
	}

	messages, err := loadMessages()
	if err != nil {
		fmt.Fprintln(logs, "Can't load messages")
		fmt.Fprintln(logs, err.Error())
		return
	}

	for _, msg := range messages {
		messagesSelection.AddItem(msg.Name, msg.Subject, 'a', func() {
			fmt.Fprintln(logs, "Selected "+msg.Name)
			activeMessage = msg
			fmt.Fprintln(messagePreview, activeMessage.Body)
		})
	}


	fmt.Println("Connecting to '" + activeConnection.Name + "'")
	client := GetClient(activeConnection)

	SendMessage(activeConnection.Destination, activeMessage, client)
    */
}
