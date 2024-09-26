package ui

import "github.com/rivo/tview"

// UI implement terminal user interface features.
type UI struct {

	// View components.
	App          *tview.Application
	Grid         *tview.Grid
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
	ui.Connections.SetTitle("Connections: ").SetBorder(true)
	ui.Messages.SetTitle("Messages: ").SetBorder(true)
    ui.Destinations.SetTitle("Destinations: ").SetBorder(true)

	//Set layouts
	navigate := tview.NewGrid().SetRows(0, 0, 0).
		AddItem(ui.Connections, 0, 0, 1, 1, 0, 0, true).
		AddItem(ui.Messages, 1, 0, 1, 1, 0, 0, false).
        AddItem(ui.Destinations, 2, 0, 1, 1, 0, 0, false)

	ui.Grid = tview.NewGrid().
		SetRows(0, 2).
		SetColumns(40, 0).
		SetBorders(false).
		AddItem(navigate, 0, 0, 1, 1, 0, 0, true)

	return &ui
}

func (ui *UI) Start() error {
	return ui.App.SetRoot(ui.Grid, true).EnableMouse(true).Run()
}
