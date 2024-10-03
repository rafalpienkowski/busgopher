package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/rafalpienkowski/busgopher/internal/controller"
)

// UI implement terminal user interface features.
type UI struct {
	controller *controller.Controller

	// View components
	theme        Theme
	App          *tview.Application
	Flex         *tview.Flex
	Connections  *tview.List
	Destinations *tview.List
	Messages     *tview.List
	Content      *tview.TextView
	Logs         *tview.TextView
	Send         *BoxButton
	Close        *BoxButton

	inputs []tview.Primitive
}

func NewUI(controller *controller.Controller) *UI {
	ui := UI{}

	ui.controller = controller

	// Create UI elements
	ui.theme = Dark()
	ui.App = tview.NewApplication()
	ui.Connections = tview.NewList().
		ShowSecondaryText(false).
		SetWrapAround(true).
		SetHighlightFullLine(true)
	ui.Destinations = tview.NewList().
		ShowSecondaryText(false).
		SetWrapAround(true).
		SetHighlightFullLine(true)
	ui.Messages = tview.NewList().
		ShowSecondaryText(false).
		SetWrapAround(true).
		SetHighlightFullLine(true)
	ui.Content = tview.NewTextView()
	ui.Logs = tview.NewTextView()
	ui.Send = ui.Send.NewBoxButton("Send").SetSelectedFunc(func() {
		destination := controller.SelectedConnection.Destinations[ui.Destinations.GetCurrentItem()]
		ui.printLog("Sending message to: " + destination)

		err := ui.controller.Send(destination)
		if err != nil {
			ui.printLog("[Error] " + err.Error())
		}
		ui.printLog("Message send")
	})
	ui.Close = ui.Close.NewBoxButton("Close").SetSelectedFunc(func() {
		ui.App.Stop()
	})

	ui.inputs = []tview.Primitive{
		ui.Connections,
		ui.Destinations,
		ui.Messages,
		ui.Content,
		ui.Send,
		ui.Close,
	}

	// Configure appearence
	ui.Connections.SetTitle(" Connections: ").SetBorder(true)
	ui.Connections.Box.SetBackgroundColor(ui.theme.backgroundColor)
	ui.Connections.SetMainTextStyle(ui.theme.style)

	ui.Destinations.SetTitle(" Destinations: ").SetBorder(true)
	ui.Destinations.Box.SetBackgroundColor(ui.theme.backgroundColor)
	ui.Destinations.SetMainTextStyle(ui.theme.style)

	ui.Messages.SetTitle(" Messages: ").SetBorder(true)
	ui.Messages.Box.SetBackgroundColor(ui.theme.backgroundColor)
	ui.Messages.SetMainTextStyle(ui.theme.style)

	ui.Content.SetTitle(" Content: ").SetBorder(true)
	ui.Content.SetBackgroundColor(ui.theme.backgroundColor)

	ui.Logs.SetTitle(" Logs: ").SetBorder(true)
	ui.Logs.SetBackgroundColor(ui.theme.backgroundColor)

	// Set layouts
	left := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.Connections, 0, 1, true).
		AddItem(ui.Destinations, 0, 1, false).
		AddItem(ui.Messages, 0, 1, false)

	actions := tview.NewFlex()
	actions.
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorGray), 0, 1, false).
		AddItem(ui.Send, ui.Send.GetWidth(), 0, false).
		AddItem(ui.Close, ui.Close.GetWidth(), 0, false)

	right := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.Content, 0, 3, false).
		AddItem(actions, 3, 0, false).
		AddItem(ui.Logs, 0, 1, false)

	ui.Flex = tview.NewFlex().
		AddItem(left, 0, 1, false).
		AddItem(right, 0, 3, false)

	ui.App.SetAfterDrawFunc(ui.setAfterDrawFunc)

	ui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			ui.cycleFocus(false)
		} else if event.Key() == tcell.KeyBacktab {
			ui.cycleFocus(true)
		}
		return event
	})

	return &ui
}

func (ui *UI) refreshDestinations() {
	ui.Destinations.Clear()
	for _, entity := range ui.controller.SelectedConnection.Destinations {
		ui.Destinations.AddItem(entity, entity, 0, func() {
			ui.printLog("Selected destination: " + entity)
		})
	}
}

func (ui *UI) LoadData() {
	for _, conn := range ui.controller.Connections {
		ui.Connections.AddItem(conn.Name, conn.Namespace, 0, func() {
			ui.controller.SelectedConnection = conn
			ui.refreshDestinations()
			ui.printLog("Selected connection: " + conn.Name + " (" + conn.Namespace + ")")
		})
	}

	for _, msg := range ui.controller.Messages {
		ui.Messages.AddItem(msg.Name, msg.Subject, 0, func() {
			ui.controller.SelectMessage(msg)
			ui.printContent(msg.Print())
			ui.printLog("Selected message: " + msg.Name)
		})
	}
}

func (ui *UI) Start() error {
	return ui.App.SetRoot(ui.Flex, true).SetFocus(ui.Connections).EnableMouse(false).Run()
}

func (ui *UI) printLog(logMsg string) {
	fmt.Fprintf(ui.Logs, "[%v]: %v\n", time.Now().Format("2006-01-02 15:04:05"), logMsg)

	getAvailableRows := func() int {
		_, _, _, height := ui.Logs.GetRect()

		return height - 2 // Minus border
	}

	ui.Logs.SetMaxLines(getAvailableRows())
}

func (ui *UI) printContent(content string) {
	ui.Content.Clear()
	fmt.Fprintf(ui.Content, "%v", content)
}

// Changes focus on TAB pressed
func (ui *UI) cycleFocus(
	reverse bool,
) {
	for i, el := range ui.inputs {
		if !el.HasFocus() {
			continue
		}

		if reverse {
			i = i - 1
			if i < 0 {
				i = len(ui.inputs) - 1
			}
		} else {
			i = i + 1
			i = i % len(ui.inputs)
		}
		ui.App.SetFocus(ui.inputs[i])
		return
	}
}

func (ui *UI) queueUpdateDraw(f func()) {
	go func() {
		ui.App.QueueUpdateDraw(f)
	}()
}

func (ui *UI) setAfterDrawFunc(screen tcell.Screen) {
	ui.queueUpdateDraw(func() {
		p := ui.App.GetFocus()

		ui.Connections.SetBorderColor(tcell.ColorWhite)
		ui.Destinations.SetBorderColor(tcell.ColorWhite)
		ui.Messages.SetBorderColor(tcell.ColorWhite)
		ui.Content.SetBorderColor(tcell.ColorWhite)
		ui.Logs.SetBorderColor(tcell.ColorWhite)
		ui.Send.SetBorderColor(tcell.ColorWhite)
		ui.Close.SetBorderColor(tcell.ColorWhite)

		switch p {
		case ui.Connections:
			ui.Connections.SetBorderColor(tcell.ColorBlue)
		case ui.Destinations:
			ui.Destinations.SetBorderColor(tcell.ColorBlue)
		case ui.Messages:
			ui.Messages.SetBorderColor(tcell.ColorBlue)
		case ui.Content:
			ui.Content.SetBorderColor(tcell.ColorBlue)
		case ui.Logs:
			ui.Logs.SetBorderColor(tcell.ColorBlue)
		case ui.Send:
			ui.Send.SetBorderColor(tcell.ColorBlue)
		case ui.Close:
			ui.Close.SetBorderColor(tcell.ColorBlue)
		}
	})
}
