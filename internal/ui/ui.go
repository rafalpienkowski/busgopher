package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI implement terminal user interface features.
type UI struct {

	// View components.
	App          *tview.Application
	Flex         *tview.Flex
	Connections  *tview.List
	Messages     *tview.List
	Destinations *tview.List
	Content      *tview.TextView
	Logs         *tview.TextView
}

// Changes focus on TAB pressed
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
		fmt.Fprintf(logs, "Focus on: %v\n", i)
		app.SetFocus(elements[i])
		return
	}
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

	inputs := []tview.Primitive{
		ui.Connections,
		ui.Messages,
		ui.Destinations,
		ui.Content,
		ui.Logs,
	}

	ui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		fmt.Fprintln(ui.Logs, "Focus changed")
		if event.Key() == tcell.KeyTab {
			cycleFocus(ui.App, inputs, false, ui.Logs)
		} else if event.Key() == tcell.KeyBacktab {
			cycleFocus(ui.App, inputs, true, ui.Logs)
		}
		return event
	})

	return &ui
}

func (ui *UI) Start() error {
	return ui.App.SetRoot(ui.Flex, true).SetFocus(ui.Connections).EnableMouse(true).Run()
}
