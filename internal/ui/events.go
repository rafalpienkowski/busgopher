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
		currentPage, _ := ui.pages.GetFrontPage()
		focusedElement := ui.app.GetFocus()
		switch currentPage {
		case "sending":
			ui.sending.setAfterDrawFunc(focusedElement)
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

func (ui *UI) cycleFocus(reverse bool) {
	currentPage, _ := ui.pages.GetFrontPage()
	if currentPage == "sending" {
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
