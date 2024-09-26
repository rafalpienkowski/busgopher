package main

import (
	"fmt"
	"os"

	"github.com/rafalpienkowski/busgopher/internal/ui"
)

func main() {

    ui := ui.NewUI()
    err := ui.Start()
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
