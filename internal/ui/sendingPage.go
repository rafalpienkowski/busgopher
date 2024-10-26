package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/rafalpienkowski/busgopher/internal/controller"
)

type SendingPage struct {
	theme      Theme
	controller *controller.Controller
	closeApp   closeAppFunc
	switchPage switchPageFunc

	flex         *tview.Flex
	connections  *tview.List
	destinations *tview.List
	messages     *tview.List
	content      *tview.TextView
	logs         *tview.TextView
	config       *BoxButton
	send         *BoxButton
	close        *BoxButton

	inputs []tview.Primitive
}

func newSendingPage(
	theme Theme,
	closeApp closeAppFunc,
	switchPage switchPageFunc,
) *SendingPage {

	flex := tview.NewFlex()
	connections := tview.NewList()
	destinations := tview.NewList()
	messages := tview.NewList()
	content := tview.NewTextView()
	logs := tview.NewTextView()
	send := newBoxButton("Send")
	config := newBoxButton("Configuration")
	close := newBoxButton("Close")

	inputs := []tview.Primitive{
		connections,
		destinations,
		messages,
		content,
		send,
		config,
		close,
	}

	sendingPage := SendingPage{
		theme:        theme,
		switchPage:   switchPage,
		closeApp:     closeApp,
		flex:         flex,
		connections:  connections,
		destinations: destinations,
		messages:     messages,
		content:      content,
		logs:         logs,
		send:         send,
		config:       config,
		close:        close,
		inputs:       inputs,
	}
	sendingPage.configureAppearence()
	sendingPage.setLayout()

	return &sendingPage
}

func (sendingPage *SendingPage) configureAppearence() {

	sendingPage.connections.
		ShowSecondaryText(false).
		SetWrapAround(true).
		SetHighlightFullLine(true).
		SetTitle(" Connections: ").
		SetBorder(true).
		SetBackgroundColor(sendingPage.theme.backgroundColor)
	sendingPage.connections.SetMainTextStyle(sendingPage.theme.style)

	sendingPage.destinations.
		ShowSecondaryText(false).
		SetWrapAround(true).
		SetHighlightFullLine(true).
		SetTitle(" Destinations: ").
		SetBorder(true).
		SetBackgroundColor(sendingPage.theme.backgroundColor)
	sendingPage.destinations.SetMainTextStyle(sendingPage.theme.style)

	sendingPage.messages.
		ShowSecondaryText(false).
		SetWrapAround(true).
		SetHighlightFullLine(true).
		SetTitle(" Messages: ").
		SetBorder(true).
		SetBackgroundColor(sendingPage.theme.backgroundColor)
	sendingPage.messages.SetMainTextStyle(sendingPage.theme.style)

	sendingPage.content.
		SetTitle(" Content: ").
		SetBorder(true).
		SetBackgroundColor(sendingPage.theme.backgroundColor)

    sendingPage.content.SetDynamicColors(true)

	sendingPage.logs.
		SetTitle(" Logs: ").
		SetBorder(true).
		SetBackgroundColor(sendingPage.theme.backgroundColor)

	sendingPage.flex.
		SetBorder(true).
		SetBackgroundColor(sendingPage.theme.backgroundColor).
		SetTitle("Sending messages")
}

func (sendingPage *SendingPage) setLayout() {
	left := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(sendingPage.connections, 0, 1, true).
		AddItem(sendingPage.destinations, 0, 1, false).
		AddItem(sendingPage.messages, 0, 1, false)

	actions := tview.NewFlex()
	actions.
		AddItem(tview.NewBox().SetBackgroundColor(sendingPage.theme.backgroundColor), 0, 1, false).
		AddItem(sendingPage.send, sendingPage.send.GetWidth(), 0, false).
		AddItem(sendingPage.config, sendingPage.config.GetWidth(), 0, false).
		AddItem(sendingPage.close, sendingPage.close.GetWidth(), 0, false)

	right := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(sendingPage.content, 0, 3, false).
		AddItem(actions, 3, 0, false).
		AddItem(sendingPage.logs, 0, 1, false)

	sendingPage.flex.
		AddItem(left, 0, 1, false).
		AddItem(right, 0, 3, false)
}

func (sendingPage *SendingPage) loadData(controller *controller.Controller) {
	sendingPage.controller = controller
	sendingPage.setActions()
	sendingPage.refreshConnections()
	sendingPage.refreshMessages()
}

func (sendingPage *SendingPage) setActions() {
	sendingPage.send.SetSelectedFunc(func() {
		err := sendingPage.controller.Send()
		if err != nil {
			sendingPage.printError(err)
		}
	})
    sendingPage.config.SetSelectedFunc(func (){
        sendingPage.switchPage("config")
    })
	sendingPage.close.SetSelectedFunc(func() {
		sendingPage.closeApp()
	})
}

func (sendingPage *SendingPage) refreshDestinations() {
	sendingPage.destinations.Clear()
	for _, name := range sendingPage.controller.GetDestiationNamesForSelectedConnection() {
		sendingPage.destinations.AddItem(name, name, 0, func() {
			err := sendingPage.controller.SelectDestinationByName(name)
			if err != nil {
				sendingPage.printError(err)
			}
		})
	}
}

func (sendingPage *SendingPage) refreshConnections() {

	sendingPage.connections.Clear()
	for _, conn := range sendingPage.controller.GetConnections() {
		sendingPage.connections.AddItem(conn.Name, conn.Namespace, 0, func() {
			err := sendingPage.controller.SelectConnectionByName(conn.Name)
			if err != nil {
				sendingPage.printError(err)
				return
			}
			sendingPage.refreshDestinations()
		})
	}
}

func (sendingPage *SendingPage) refreshMessages() {
	sendingPage.messages.Clear()

	for _, msg := range sendingPage.controller.GetMessages() {
		sendingPage.messages.AddItem(msg.Name, msg.Message.Subject, 0, func() {
			err := sendingPage.controller.SelectMessageByName(msg.Name)
			if err != nil {
				sendingPage.printError(err)
				return
			}
            colorized, err := colorizeJSON(msg.Message.Print())
            if err != nil {
                sendingPage.printError(err)
            }
			sendingPage.printContent(colorized)
		})
	}
}

func (sendingPage *SendingPage) printContent(content string) {
	sendingPage.content.Clear()
	fmt.Fprintf(sendingPage.content, "%v", content)
}

func (sendingPage *SendingPage) printError(err error) {
	sendingPage.printLog(fmt.Sprintf(
		"[%v]: [Error] %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		err.Error(),
	))
}

func (sendingPage *SendingPage) printLog(logMsg string) {
	fmt.Fprintf(sendingPage.logs, "%v", logMsg)

	getAvailableRows := func() int {
		_, _, _, height := sendingPage.logs.GetRect()

		return height - 2 // Minus border
	}

	sendingPage.logs.SetMaxLines(getAvailableRows())
}

func (sendingPage *SendingPage) setAfterDrawFunc(focusedElement tview.Primitive) {
	sendingPage.connections.SetBorderColor(tcell.ColorWhite)
	sendingPage.destinations.SetBorderColor(tcell.ColorWhite)
	sendingPage.messages.SetBorderColor(tcell.ColorWhite)
	sendingPage.content.SetBorderColor(tcell.ColorWhite)
	sendingPage.logs.SetBorderColor(tcell.ColorWhite)
	sendingPage.send.SetBorderColor(tcell.ColorWhite)
	sendingPage.config.SetBorderColor(tcell.ColorWhite)
	sendingPage.close.SetBorderColor(tcell.ColorWhite)

	switch focusedElement {
	case sendingPage.connections:
		sendingPage.connections.SetBorderColor(tcell.ColorBlue)
	case sendingPage.destinations:
		sendingPage.destinations.SetBorderColor(tcell.ColorBlue)
	case sendingPage.messages:
		sendingPage.messages.SetBorderColor(tcell.ColorBlue)
	case sendingPage.content:
		sendingPage.content.SetBorderColor(tcell.ColorBlue)
	case sendingPage.logs:
		sendingPage.logs.SetBorderColor(tcell.ColorBlue)
	case sendingPage.send:
		sendingPage.send.SetBorderColor(tcell.ColorBlue)
	case sendingPage.config:
		sendingPage.config.SetBorderColor(tcell.ColorBlue)
	case sendingPage.close:
		sendingPage.close.SetBorderColor(tcell.ColorBlue)
	}
}

func colorizeJSON(data string) (string, error) {
	var jsonData interface{}
	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		return "", err
	}
	var buf bytes.Buffer
	formatJSON(jsonData, &buf, "")
	return buf.String(), nil
}

func formatJSON(value interface{}, buf *bytes.Buffer, indent string) {
	switch v := value.(type) {
	case map[string]interface{}:
		buf.WriteString("{\n")
		for key, val := range v {
			buf.WriteString(indent + "  ")
			buf.WriteString(fmt.Sprintf(`[yellow]"%s":[-] `, key))
			formatJSON(val, buf, indent+"  ")
			buf.WriteString(",\n")
		}
		buf.WriteString(indent + "}")
	case []interface{}:
		buf.WriteString("[\n")
		for _, val := range v {
			buf.WriteString(indent + "  ")
			formatJSON(val, buf, indent+"  ")
			buf.WriteString(",\n")
		}
		buf.WriteString(indent + "]")
	case string:
		buf.WriteString(fmt.Sprintf(`[green]"%s"[-]`, v))
	case float64:
		buf.WriteString(fmt.Sprintf(`[cyan]%v[-]`, v))
	case bool:
		buf.WriteString(fmt.Sprintf(`[magenta]%v[-]`, v))
	case nil:
		buf.WriteString(`[gray]null[-]`)
	default:
		buf.WriteString(fmt.Sprintf(`%v`, v))
	}
}
