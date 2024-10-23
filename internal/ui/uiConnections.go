package ui

import (
	"maps"
	"slices"

	"github.com/rivo/tview"

	"github.com/rafalpienkowski/busgopher/internal/asb"
)

func (ui *UI) addConnection() {
	ui.PrintLog("Addding new connection\n")

	ui.advancedForm.clear()
	ui.advancedForm.form.
		AddInputField("Connection name", "", 0, nil, nil).
		AddInputField("Service Bus namespace", "", 0, nil, nil).
		AddButton("Add", func() {
			ui.advancedForm.message.Clear()
			newConnection := asb.Connection{
				Name:         ui.advancedForm.form.GetFormItem(0).(*tview.InputField).GetText(),
				Namespace:    ui.advancedForm.form.GetFormItem(1).(*tview.InputField).GetText(),
				Destinations: []string{},
			}

			err := ui.controller.AddConnection(&newConnection)
			if err != nil {
				ui.advancedForm.setMessage(err.Error())
				return
			}

			ui.refreshConnections()
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		}).
		AddButton("Quit", func() {
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		})

	ui.switchToForm("Add new connection")
}

func (ui *UI) removeConnection() {
	ui.PrintLog("Removing connection\n")
	ui.advancedForm.clear()
	ui.advancedForm.setMessage("Please select connection to remove")
	ui.advancedForm.form.
		AddDropDown(
			"Connection to remove",
			slices.Collect(maps.Keys(ui.controller.Config.NConnections)),
			0,
			nil,
		).
		AddButton("Remove", func() {

			idx, option := ui.advancedForm.form.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
			if idx < 0 {
				ui.advancedForm.setMessage("Please select connection to remove")
				return
			}

			err := ui.controller.RemoveConnection(option)
			if err != nil {
				ui.advancedForm.setMessage(err.Error())
				return
			}
            ui.PrintLog("Removed connection: " + option)
			ui.refreshConnections()
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		}).
		AddButton("Quit", func() {
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		})

	ui.switchToForm("Remove connection")
}
