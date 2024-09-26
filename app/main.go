package main

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var activeConnection busConnection
var activeMessage busMessage

func cycleFocus(
	app *tview.Application,
	elements []tview.Primitive,
	reverse bool,
	logs *tview.TextView,
) {
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
		fmt.Fprintln(logs, "Focus on: "+strconv.Itoa(i))
		app.SetFocus(elements[i])
		return
	}
}

func main() {

	app := tview.NewApplication()
	logs := tview.NewTextView()
    logs.SetBorder(true)
	connectionsSelection := tview.NewList()
    connectionsSelection.SetBorder(true).SetTitle("Connection")
	messagesSelection := tview.NewList()
    messagesSelection.SetBorder(true).SetTitle("Message")
	messagePreview := tview.NewTextView()
    messagePreview.SetBorder(true)
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

	messagePreview.
		SetDynamicColors(true).
		SetWrap(true).
		SetBorder(true).
		SetTitle("Message")

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

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
            AddItem(connectionsSelection, 0, 1, true).
			AddItem(messagesSelection, 0, 1, false),
			0, 2, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Right "), 0, 5, false).
			AddItem(logs, 10, 1, false),
			0, 5, false)

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
