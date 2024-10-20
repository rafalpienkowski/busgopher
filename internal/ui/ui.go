package ui

import (
	"fmt"

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

func NewUI() *UI {
	ui := UI{}

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
		ui.controller.Send()
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
	ui.Connections.SetBackgroundColor(ui.theme.backgroundColor)
	ui.Connections.SetMainTextStyle(ui.theme.style)

	ui.Destinations.SetTitle(" Destinations: ").SetBorder(true)
	ui.Destinations.SetBackgroundColor(ui.theme.backgroundColor)
	ui.Destinations.SetMainTextStyle(ui.theme.style)

	ui.Messages.SetTitle(" Messages: ").SetBorder(true)
	ui.Messages.SetBackgroundColor(ui.theme.backgroundColor)
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
	ui.App.SetInputCapture(ui.setInputCapture)

	return &ui
}

func (ui *UI) refreshDestinations() {
	ui.Destinations.Clear()
	for _, name := range ui.controller.GetDestiationNamesForSelectedConnection() {
		ui.Destinations.AddItem(name, name, 0, func() {
			ui.controller.SelectDestinationByName(name)
		})
	}
}

func (ui *UI) LoadData(controller *controller.Controller) {

	ui.controller = controller

	for _, conn := range ui.controller.Config.NConnections {
		ui.Connections.AddItem(conn.Name, conn.Namespace, 0, func() {
			ui.controller.SelectConnectionByName(conn.Name)
			ui.refreshDestinations()
		})
	}

	for _, msg := range ui.controller.Config.NMessages {
		ui.Messages.AddItem(msg.Name, msg.Subject, 0, func() {
			ui.controller.SelectMessageByName(msg.Name)
			ui.printContent(msg.Print())
		})
	}
}

func (ui *UI) Start() error {
	return ui.App.SetRoot(ui.Flex, true).SetFocus(ui.Connections).EnableMouse(false).Run()
}

func (ui *UI) PrintLog(logMsg string) {
	fmt.Fprintf(ui.Logs, "%v", logMsg)

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
