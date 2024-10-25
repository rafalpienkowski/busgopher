package ui

import (
	"github.com/rivo/tview"

	"github.com/rafalpienkowski/busgopher/internal/controller"
)

// UI implement terminal user interface features.
type UI struct {
	controller *controller.Controller

	// View components
	theme Theme
	app   *tview.Application

	// Pages
	pages   *tview.Pages
	sending *SendingPage

	// Config
	configFlex   *tview.Flex
	advancedForm *AdvancedForm

	inputs []tview.Primitive
}

type closeAppFunc func()

func NewUI() *UI {
	ui := UI{}

	// Create UI elements
	ui.theme = Dark()
	ui.app = tview.NewApplication()
	ui.pages = tview.NewPages()
	ui.sending = newSendingPage(ui.theme, closeAppFunc(func() { ui.app.Stop() }))

	ui.configFlex = tview.NewFlex()
	ui.advancedForm = newAdvancedForm(ui.theme)

	ui.inputs = []tview.Primitive{
		ui.sending.connections,
		ui.sending.destinations,
		ui.sending.messages,
		ui.sending.content,
		ui.sending.send,
		ui.sending.close,
	}

	ui.pages.
		AddPage("sending", ui.sending.flex, true, true).
		AddPage("form", ui.advancedForm.flex, true, false)

	ui.app.SetAfterDrawFunc(ui.setAfterDrawFunc)
	ui.app.SetInputCapture(ui.setInputCapture)

	return &ui
}

func (ui *UI) LoadData(controller *controller.Controller) {
	ui.controller = controller
    ui.sending.loadData(ui.controller)
}

func (ui *UI) Start() error {
	return ui.app.SetRoot(ui.pages, true).
		SetFocus(ui.sending.connections).
		EnableMouse(false).
		Run()
}

func (ui *UI) switchToPage(title string, page string) {
	ui.advancedForm.flex.SetTitle(title)
	ui.pages.SwitchToPage(page)
	ui.app.SetFocus(ui.advancedForm.flex)
	ui.app.SetFocus(ui.advancedForm.form)
}
