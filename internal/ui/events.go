package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (ui *UI) queueUpdateDraw(f func()) {
	go func() {
		ui.App.QueueUpdateDraw(f)
	}()
}

func (ui *UI) setAfterDrawFunc(screen tcell.Screen) {
	ui.queueUpdateDraw(func() {
		p := ui.App.GetFocus()

		ui.connections.SetBorderColor(tcell.ColorWhite)
		ui.destinations.SetBorderColor(tcell.ColorWhite)
		ui.messages.SetBorderColor(tcell.ColorWhite)
		ui.content.SetBorderColor(tcell.ColorWhite)
		ui.Logs.SetBorderColor(tcell.ColorWhite)
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
		case ui.Logs:
			ui.Logs.SetBorderColor(tcell.ColorBlue)
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
		p := ui.App.GetFocus()
		switch p {
		case ui.connections:
			ui.addConnection()
		case ui.destinations:
			ui.PrintLog("New destination")
		case ui.messages:
			ui.PrintLog("New message")
		}
	case tcell.KeyCtrlU:
		p := ui.App.GetFocus()
		switch p {
		case ui.connections:
			ui.PrintLog("Update connection")
		case ui.destinations:
			ui.PrintLog("Update destination")
		case ui.messages:
			ui.PrintLog("Update message")
		}
	case tcell.KeyCtrlD:
		p := ui.App.GetFocus()
		switch p {
		case ui.connections:
			ui.PrintLog("Delete connection")
		case ui.destinations:
			ui.PrintLog("Delete destination")
		case ui.messages:
			ui.PrintLog("Delete message")
		}
	}

	return event
}

func (ui *UI) addConnection() {
	ui.PrintLog("Addding new connection\n")

	ui.form = ui.form.Clear(true)
	ui.form.
		AddInputField("Connection name", "", 0, nil, nil).
        AddInputField("Service Bus namespace", "", 0, nil, nil).
		AddButton("Add", func() {
			ui.pages.SwitchToPage("sending")
			ui.App.SetFocus(ui.connections)
		}).
		AddButton("Quit", func() {
			ui.pages.SwitchToPage("sending")
			ui.App.SetFocus(ui.connections)
		})

	for idx := range (ui.form.GetButtonCount() - 1) {
		ui.form.GetFormItem(idx).(*tview.InputField).
            SetFieldBackgroundColor(ui.theme.backgroundColor).
			SetFieldTextColor(ui.theme.backgroundColor).
			SetLabelColor(ui.theme.foregroundColor)
	}

	ui.formFlex.SetTitle("Add new connection")
	ui.pages.SwitchToPage("form")
	ui.App.SetFocus(ui.formFlex)
	ui.App.SetFocus(ui.form)
}

func (ui *UI) cycleFocus(reverse bool) {
	currentPage, _ := ui.pages.GetFrontPage()
	if currentPage == "form" {
		ui.PrintLog("form next")
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
		ui.App.SetFocus(ui.inputs[i])
		return
	}
}
