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
	app          *tview.Application
	pages        *tview.Pages
	sendingFlex  *tview.Flex
	connections  *tview.List
	destinations *tview.List
	messages     *tview.List
	content      *tview.TextView
	logs         *tview.TextView
	send         *BoxButton
	close        *BoxButton
	advancedForm *AdvancedForm

	inputs []tview.Primitive
}

func NewUI() *UI {
	ui := UI{}

	// Create UI elements
	ui.theme = Dark()
	ui.app = tview.NewApplication()
	ui.pages = tview.NewPages()
	ui.sendingFlex = tview.NewFlex()
	ui.advancedForm = newAdvancedForm(ui.theme)

	ui.connections = tview.NewList().
		ShowSecondaryText(false).
		SetWrapAround(true).
		SetHighlightFullLine(true)
	ui.destinations = tview.NewList().
		ShowSecondaryText(false).
		SetWrapAround(true).
		SetHighlightFullLine(true)
	ui.messages = tview.NewList().
		ShowSecondaryText(false).
		SetWrapAround(true).
		SetHighlightFullLine(true)
	ui.content = tview.NewTextView()
	ui.logs = tview.NewTextView()
	ui.send = ui.send.NewBoxButton("Send").SetSelectedFunc(func() {
		ui.controller.Send()
	})
	ui.close = ui.close.NewBoxButton("Close").SetSelectedFunc(func() {
		ui.app.Stop()
	})

	ui.inputs = []tview.Primitive{
		ui.connections,
		ui.destinations,
		ui.messages,
		ui.content,
		ui.send,
		ui.close,
	}

	// Configure appearence
	ui.connections.SetTitle(" Connections: ").SetBorder(true)
	ui.connections.SetBackgroundColor(ui.theme.backgroundColor)
	ui.connections.SetMainTextStyle(ui.theme.style)

	ui.destinations.SetTitle(" Destinations: ").SetBorder(true)
	ui.destinations.SetBackgroundColor(ui.theme.backgroundColor)
	ui.destinations.SetMainTextStyle(ui.theme.style)

	ui.messages.SetTitle(" Messages: ").SetBorder(true)
	ui.messages.SetBackgroundColor(ui.theme.backgroundColor)
	ui.messages.SetMainTextStyle(ui.theme.style)

	ui.content.SetTitle(" Content: ").SetBorder(true)
	ui.content.SetBackgroundColor(ui.theme.backgroundColor)

	ui.logs.SetTitle(" Logs: ").SetBorder(true)
	ui.logs.SetBackgroundColor(ui.theme.backgroundColor)

	// Set layouts
	left := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.connections, 0, 1, true).
		AddItem(ui.destinations, 0, 1, false).
		AddItem(ui.messages, 0, 1, false)

	actions := tview.NewFlex()
	actions.
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorGray), 0, 1, false).
		AddItem(ui.send, ui.send.GetWidth(), 0, false).
		AddItem(ui.close, ui.close.GetWidth(), 0, false)

	right := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.content, 0, 3, false).
		AddItem(actions, 3, 0, false).
		AddItem(ui.logs, 0, 1, false)

	ui.sendingFlex.
		SetBorder(true).
		SetBackgroundColor(ui.theme.backgroundColor).
		SetTitle("Sending messages")

	ui.sendingFlex.
		AddItem(left, 0, 1, false).
		AddItem(right, 0, 3, false)

	ui.pages.
		AddPage("sending", ui.sendingFlex, true, true).
		AddPage("form", ui.advancedForm.flex, true, false)

	ui.app.SetAfterDrawFunc(ui.setAfterDrawFunc)
	ui.app.SetInputCapture(ui.setInputCapture)

	return &ui
}

func (ui *UI) refreshDestinations() {
	ui.destinations.Clear()
	for _, name := range ui.controller.GetDestiationNamesForSelectedConnection() {
		ui.destinations.AddItem(name, name, 0, func() {
			ui.controller.SelectDestinationByName(name)
		})
	}
}

func (ui *UI) LoadData(controller *controller.Controller) {
	ui.controller = controller
	ui.refreshConnections()
	ui.refreshMessages()
}

func (ui *UI) refreshConnections() {

	ui.connections.Clear()
	for _, conn := range ui.controller.Config.NConnections {
		ui.connections.AddItem(conn.Name, conn.Namespace, 0, func() {
			ui.controller.SelectConnectionByName(conn.Name)
			ui.refreshDestinations()
		})
	}
}

func (ui *UI) refreshMessages() {
	ui.messages.Clear()

	for _, msg := range ui.controller.Config.NMessages {
		ui.messages.AddItem(msg.Name, msg.Subject, 0, func() {
			ui.controller.SelectMessageByName(msg.Name)
			ui.printContent(msg.Print())
		})
	}
}

func (ui *UI) Start() error {
	return ui.app.SetRoot(ui.pages, true).SetFocus(ui.connections).EnableMouse(true).Run()
}

func (ui *UI) PrintLog(logMsg string) {
	fmt.Fprintf(ui.logs, "%v", logMsg)

	getAvailableRows := func() int {
		_, _, _, height := ui.logs.GetRect()

		return height - 2 // Minus border
	}

	ui.logs.SetMaxLines(getAvailableRows())
}

func (ui *UI) printContent(content string) {
	ui.content.Clear()
	fmt.Fprintf(ui.content, "%v", content)
}

func (ui *UI) switchToForm(title string){
	ui.advancedForm.flex.SetTitle(title)
	ui.pages.SwitchToPage("form")
	ui.app.SetFocus(ui.advancedForm.flex)
	ui.app.SetFocus(ui.advancedForm.form)
}
