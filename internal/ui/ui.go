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
}

func NewUI() *UI {
	ui := UI{}

	// Create UI elements
	ui.App = tview.NewApplication()
	ui.Connections = tview.NewList().ShowSecondaryText(false)
	ui.Messages = tview.NewList().ShowSecondaryText(false)
	ui.Destinations = tview.NewList().ShowSecondaryText(false)

	// Configure appearence
	ui.Connections.SetTitle(" Connections: ").SetBorder(true)
	ui.Messages.SetTitle(" Messages: ").SetBorder(true)
	ui.Destinations.SetTitle(" Destinations: ").SetBorder(true)

	//Set layouts
	navigate := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.Connections, 0, 1, true).
		AddItem(ui.Messages, 0, 1, false).
		AddItem(ui.Destinations, 0, 1, false)

	ui.Flex = tview.NewFlex().
		AddItem(navigate, 0, 1, false)

	return &ui
}

func (ui *UI) Start() error {
	return ui.App.SetRoot(ui.Flex, true).EnableMouse(true).Run()
}
