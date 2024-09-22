package main

import (
	"fmt"

	"github.com/rivo/tview"
)

var activeConnection busConnection
var activeMessage busMessage

func main() {

	app := tview.NewApplication()

	logs := tview.NewTextView()

	logs.
		SetDynamicColors(true).
		SetWrap(true).
		SetBorder(true).
		SetTitle("Logs")

	connections, err := loadConnections()
	if err != nil {
		fmt.Fprintln(logs, "Can't load connections")
		fmt.Fprintln(logs, err.Error())
		return
	}

    var connectionNames []string
    for _, conn := range connections {
        connectionNames = append(connectionNames, conn.Name)
    }
    connectionNames = append(connectionNames, "Add..")

	connectionsDropDown := tview.NewDropDown().
		SetOptions(connectionNames, func(name string, index int) {
			fmt.Fprintln(logs, "Selected "+name)
		}).
		SetLabel("Select connection: ")

	flex := tview.NewFlex().
		AddItem(connectionsDropDown, 0, 1, true).
		AddItem(logs, 0, 3, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	return
	activeConnection = connections[0]

	messages, err := loadMessages()
	if err != nil {
		fmt.Println("Can't load messages")
		fmt.Println(err.Error())
		return
	}
	activeMessage = messages[0]

	fmt.Println("Connecting to '" + activeConnection.Name + "'")
	client := GetClient(activeConnection)

	fmt.Println(
		"Sending a message '" + activeMessage.Body + "' to: '" + activeConnection.Destination + "'",
	)
	SendMessage(activeConnection.Destination, activeMessage, client)

	fmt.Println("Done")
}
