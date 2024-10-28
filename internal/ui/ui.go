package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/rafalpienkowski/busgopher/internal/controller"
)

type UI struct {
	controller *controller.Controller

	// View components
	theme Theme
	app   *tview.Application

	// Pages
	pages   *tview.Pages
	sending *SendingPage
	config  *ConfigPage
}

type closeAppFunc func()
type switchPageFunc func(string)

func NewUI() *UI {
	ui := UI{}

	ui.theme = Dark()
	ui.app = tview.NewApplication()
	ui.pages = tview.NewPages()
	ui.sending = newSendingPage(ui.theme, ui.app.Stop, ui.switchToPage)
	ui.config = newConfigPage(ui.theme, ui.app.Stop, ui.switchToPage)

	ui.pages.
		AddPage("sending", ui.sending.flex, true, true).
		AddPage("config", ui.config.flex, true, false)

	ui.app.SetAfterDrawFunc(ui.setAfterDrawFunc)
	ui.app.SetInputCapture(ui.setInputCapture)

	return &ui
}

func (ui *UI) LoadData(controller *controller.Controller) {
	ui.controller = controller
	ui.sending.loadData(ui.controller)
	ui.config.loadData(ui.controller)
}

func (ui *UI) Start() error {
	return ui.app.SetRoot(ui.pages, true).
		SetFocus(ui.sending.connections).
		EnableMouse(false).
		Run()
}

func (ui *UI) switchToPage(page string) {
	ui.pages.SwitchToPage(page)
	switch page {
	case "sending":
		ui.sending.refresh()
		ui.app.SetFocus(ui.sending.connections)
	case "config":
		ui.config.refresh()
		ui.app.SetFocus(ui.config.config)
	}
}

func (ui *UI) WriteLog(log string) {
	page, _ := ui.pages.GetFrontPage()
	switch page {
	case "sending":
		ui.sending.printLog(fmt.Sprintf(
			"[%v]: [Info] %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			log))
	case "config":
		ui.config.printLog(fmt.Sprintf(
			"[%v]: [Info] %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			log))
	}
}

func (ui *UI) queueUpdateDraw(f func()) {
	go func() {
		ui.app.QueueUpdateDraw(f)
	}()
}

func (ui *UI) setAfterDrawFunc(screen tcell.Screen) {
	ui.queueUpdateDraw(func() {
		currentPage, _ := ui.pages.GetFrontPage()
		focusedElement := ui.app.GetFocus()
		switch currentPage {
		case "sending":
			ui.sending.setAfterDrawFunc(focusedElement)
		case "config":
			ui.config.setAfterDrawFunc(focusedElement)
		}
	})
}

func (ui *UI) setInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyTab:
		ui.cycleFocus(false)
	case tcell.KeyBacktab:
		ui.cycleFocus(true)
	}

	return event
}

func (ui *UI) getNextFocusInput(inputs []tview.Primitive, reverse bool) tview.Primitive {
	for i, el := range inputs {
		if !el.HasFocus() {
			continue
		}

		if reverse {
			i = i - 1
			if i < 0 {
				i = len(inputs) - 1
			}
		} else {
			i = i + 1
			i = i % len(inputs)
		}
		return inputs[i]
	}
	return inputs[0]
}

func (ui *UI) cycleFocus(reverse bool) {
	currentPage, _ := ui.pages.GetFrontPage()

	var input tview.Primitive
	switch currentPage {
	case "sending":
		input = ui.getNextFocusInput(ui.sending.inputs, reverse)
	case "config":
		input = ui.getNextFocusInput(ui.config.inputs, reverse)
	}

	ui.app.SetFocus(input)
}
