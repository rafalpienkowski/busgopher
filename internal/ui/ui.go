package ui

import "github.com/rivo/tview"

// UI implement terminal user interface features.
type UI struct {

	// View components.
	App          *tview.Application
	Flex         *tview.Flex
	Connections  *tview.List
	Messages     *tview.List
	Destinations *tview.List
	Content         *tview.TextView
	Logs         *tview.TextView
}

func NewUI() *UI {
	ui := UI{}

	// Create UI elements
	ui.App = tview.NewApplication()
	ui.Connections = tview.NewList().ShowSecondaryText(false)
	ui.Messages = tview.NewList().ShowSecondaryText(false)
	ui.Destinations = tview.NewList().ShowSecondaryText(false)
    ui.Content = tview.NewTextView()
    ui.Logs = tview.NewTextView()

	// Configure appearence
	ui.Connections.SetTitle(" Connections: ").SetBorder(true)
	ui.Messages.SetTitle(" Messages: ").SetBorder(true)
	ui.Destinations.SetTitle(" Destinations: ").SetBorder(true)
    ui.Content.SetTitle(" Content: ").SetBorder(true)
    ui.Logs.SetTitle(" Logs: ").SetBorder(true)

	//Set layouts
	left := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.Connections, 0, 1, true).
		AddItem(ui.Messages, 0, 1, false).
		AddItem(ui.Destinations, 0, 1, false)

    right := tview.NewFlex().SetDirection(tview.FlexRow).
        AddItem(ui.Content, 0, 3, false).
        AddItem(ui.Logs, 0, 1, false)

	ui.Flex = tview.NewFlex().
		AddItem(left, 0, 1, false).
        AddItem(right, 0, 3, false)

	return &ui
}

func (ui *UI) Start() error {
	return ui.App.SetRoot(ui.Flex, true).EnableMouse(true).Run()
}
