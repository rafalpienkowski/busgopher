package ui

import "github.com/gdamore/tcell/v2"

func (ui *UI) queueUpdateDraw(f func()) {
	go func() {
		ui.App.QueueUpdateDraw(f)
	}()
}

func (ui *UI) setAfterDrawFunc(screen tcell.Screen) {
	ui.queueUpdateDraw(func() {
		p := ui.App.GetFocus()

		ui.Connections.SetBorderColor(tcell.ColorWhite)
		ui.Destinations.SetBorderColor(tcell.ColorWhite)
		ui.Messages.SetBorderColor(tcell.ColorWhite)
		ui.Content.SetBorderColor(tcell.ColorWhite)
		ui.Logs.SetBorderColor(tcell.ColorWhite)
		ui.Send.SetBorderColor(tcell.ColorWhite)
		ui.Close.SetBorderColor(tcell.ColorWhite)

		switch p {
		case ui.Connections:
			ui.Connections.SetBorderColor(tcell.ColorBlue)
		case ui.Destinations:
			ui.Destinations.SetBorderColor(tcell.ColorBlue)
		case ui.Messages:
			ui.Messages.SetBorderColor(tcell.ColorBlue)
		case ui.Content:
			ui.Content.SetBorderColor(tcell.ColorBlue)
		case ui.Logs:
			ui.Logs.SetBorderColor(tcell.ColorBlue)
		case ui.Send:
			ui.Send.SetBorderColor(tcell.ColorBlue)
		case ui.Close:
			ui.Close.SetBorderColor(tcell.ColorBlue)
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
		ui.App.SetFocus(ui.inputs[i])
		return
	}
}
