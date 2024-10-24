package ui

import (
	"github.com/gdamore/tcell/v2"
)

func (ui *UI) queueUpdateDraw(f func()) {
	go func() {
		ui.app.QueueUpdateDraw(f)
	}()
}

func (ui *UI) setAfterDrawFunc(screen tcell.Screen) {
	ui.queueUpdateDraw(func() {
		p := ui.app.GetFocus()

		ui.connections.SetBorderColor(tcell.ColorWhite)
		ui.destinations.SetBorderColor(tcell.ColorWhite)
		ui.messages.SetBorderColor(tcell.ColorWhite)
		ui.content.SetBorderColor(tcell.ColorWhite)
		ui.logs.SetBorderColor(tcell.ColorWhite)
		ui.send.SetBorderColor(tcell.ColorWhite)
		ui.close.SetBorderColor(tcell.ColorWhite)

		switch p {
		case ui.connections:
			ui.connections.SetBorderColor(tcell.ColorBlue)
		case ui.destinations:
			ui.destinations.SetBorderColor(tcell.ColorBlue)
		case ui.messages:
			ui.messages.SetBorderColor(tcell.ColorBlue)
		case ui.content:
			ui.content.SetBorderColor(tcell.ColorBlue)
		case ui.logs:
			ui.logs.SetBorderColor(tcell.ColorBlue)
		case ui.send:
			ui.send.SetBorderColor(tcell.ColorBlue)
		case ui.close:
			ui.close.SetBorderColor(tcell.ColorBlue)
		}
	})
}

func (ui *UI) setInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyTab:
		ui.cycleFocus(false)
	case tcell.KeyBacktab:
		ui.cycleFocus(true)
	case tcell.KeyCtrlN:
		p := ui.app.GetFocus()
		switch p {
		case ui.connections:
			ui.addConnection()
		case ui.destinations:
			ui.PrintLog("New destination")
		case ui.messages:
            ui.addMessage()
		}
	case tcell.KeyCtrlU:
		p := ui.app.GetFocus()
		switch p {
		case ui.connections:
			ui.updateSelectedConnection()
		case ui.destinations:
			ui.PrintLog("Update destination")
		case ui.messages:
            ui.updateSelectedMessage()
		}
	case tcell.KeyCtrlD:
		p := ui.app.GetFocus()
		switch p {
		case ui.connections:
			ui.removeSelectedConnection()
		case ui.destinations:
			ui.PrintLog("Delete destination")
		case ui.messages:
            ui.removeSelectedMessage()
		}
	}

	return event
}

func (ui *UI) cycleFocus(reverse bool) {
	currentPage, _ := ui.pages.GetFrontPage()
	if currentPage == "form" {
		return
	}

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
		ui.app.SetFocus(ui.inputs[i])
		return
	}
}
