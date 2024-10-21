package ui

import (
	"github.com/rafalpienkowski/busgopher/internal/asb"
	"github.com/rivo/tview"
)

func (ui *UI) addConnection() {
	ui.PrintLog("Addding new connection\n")

	ui.form = ui.form.Clear(true)
	ui.form.
		AddInputField("Connection name", "", 0, nil, nil).
        AddInputField("Service Bus namespace", "", 0, nil, nil).
		AddButton("Add", func() {
            newConnection := asb.Connection{
                Name: ui.form.GetFormItem(0).(*tview.InputField).GetText(),
                Namespace: ui.form.GetFormItem(1).(*tview.InputField).GetText(),
                Destinations: []string{},
            }

            err := ui.controller.AddConnection(&newConnection)
            if err != nil {
                ui.PrintLog("new error")
                ui.PrintLog(err.Error())
            }

            ui.refreshConnections()
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		}).
		AddButton("Quit", func() {
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		}).
        SetFieldBackgroundColor(ui.theme.foregroundColor).
        SetFieldTextColor(ui.theme.backgroundColor)

	ui.formFlex.SetTitle("Add new connection")
	ui.pages.SwitchToPage("form")
	ui.app.SetFocus(ui.formFlex)
	ui.app.SetFocus(ui.form)
}
