package main

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var activeConnection busConnection
var activeMessage busMessage

func cycleFocus(app *tview.Application, elements []tview.Primitive, reverse bool, logs *tview.TextView) {
	for i, el := range elements {
		if !el.HasFocus() {
			continue
		}

		if reverse {
			i = i - 1
			if i < 0 {
				i = len(elements) - 1
			}
		} else {
			i = i + 1
			i = i % len(elements)
		}
        fmt.Fprintln(logs, "Focus on: " + strconv.Itoa(i))

		app.SetFocus(elements[i])
		return
	}
}

func main() {

	app := tview.NewApplication()
	logs := tview.NewTextView()
	connectionsDropDown := tview.NewDropDown()
	messagesDropDown := tview.NewDropDown()

	inputs := []tview.Primitive{
		connectionsDropDown,
		messagesDropDown,
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			cycleFocus(app, inputs, false, logs)
		} else if event.Key() == tcell.KeyBacktab {
			cycleFocus(app, inputs, true, logs)
		}
		return event
	})

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

	connectionsDropDown.SetOptions(connectionNames, func(name string, index int) {
		fmt.Fprintln(logs, "Selected "+name)
		activeConnection = connections[index]
	}).SetLabel("Select connection: ")

	messages, err := loadMessages()
	if err != nil {
		fmt.Fprintln(logs, "Can't load messages")
		fmt.Fprintln(logs, err.Error())
		return
	}

	var messageNames []string
	for _, msg := range messages {
		messageNames = append(messageNames, msg.Name)
	}
	messageNames = append(messageNames, "Add..")

	messagesDropDown.SetOptions(messageNames, func(name string, index int) {
		fmt.Fprintln(logs, "Selected "+name)
		activeMessage = messages[index]
	}).
		SetLabel("Select message: ")

	flex := tview.NewFlex().
		AddItem(connectionsDropDown, 0, 1, true).
		AddItem(messagesDropDown, 0, 2, false).
		AddItem(logs, 0, 3, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	return

	fmt.Println("Connecting to '" + activeConnection.Name + "'")
	client := GetClient(activeConnection)

	fmt.Println(
		"Sending a message '" + activeMessage.Body + "' to: '" + activeConnection.Destination + "'",
	)
	SendMessage(activeConnection.Destination, activeMessage, client)

	fmt.Println("Done")
}
