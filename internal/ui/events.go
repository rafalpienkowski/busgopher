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
	case tcell.KeyCtrlN:
		p := ui.App.GetFocus()
		switch p {
		case ui.Connections:
			ui.addConnection()
		case ui.Destinations:
			ui.PrintLog("New destination")
		case ui.Messages:
			ui.PrintLog("New message")
		}
	case tcell.KeyCtrlU:
		p := ui.App.GetFocus()
		switch p {
		case ui.Connections:
			ui.PrintLog("Update connection")
		case ui.Destinations:
			ui.PrintLog("Update destination")
		case ui.Messages:
			ui.PrintLog("Update message")
		}
	case tcell.KeyCtrlD:
		p := ui.App.GetFocus()
		switch p {
		case ui.Connections:
			ui.PrintLog("Delete connection")
		case ui.Destinations:
			ui.PrintLog("Delete destination")
		case ui.Messages:
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
			ui.App.SetFocus(ui.Connections)
		}).
		AddButton("Quit", func() {
			ui.pages.SwitchToPage("sending")
			ui.App.SetFocus(ui.Connections)
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
