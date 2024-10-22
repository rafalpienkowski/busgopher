package ui

import (
	"github.com/rafalpienkowski/busgopher/internal/asb"
	"github.com/rivo/tview"
)

func (ui *UI) addConnection() {
	ui.PrintLog("Addding new connection\n")

	ui.advancedForm.Clear()
	ui.advancedForm.form.
		AddInputField("Connection name", "", 0, nil, nil).
        AddInputField("Service Bus namespace", "", 0, nil, nil).
		AddButton("Add", func() {
            ui.advancedForm.message.Clear()
            newConnection := asb.Connection{
                Name: ui.advancedForm.form.GetFormItem(0).(*tview.InputField).GetText(),
                Namespace: ui.advancedForm.form.GetFormItem(1).(*tview.InputField).GetText(),
                Destinations: []string{},
            }

            err := ui.controller.AddConnection(&newConnection)
            if err != nil {
                ui.advancedForm.SetMessage(err.Error())
                return
            }

            ui.refreshConnections()
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		}).
		AddButton("Quit", func() {
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		}).
        //add errors
        SetFieldBackgroundColor(ui.theme.foregroundColor).
        SetFieldTextColor(ui.theme.backgroundColor)

	ui.advancedForm.flex.SetTitle("Add new connection")
	ui.pages.SwitchToPage("form")
	ui.app.SetFocus(ui.advancedForm.flex)
	ui.app.SetFocus(ui.advancedForm.form)
}
