package ui

import (
	"github.com/rivo/tview"

	"github.com/rafalpienkowski/busgopher/internal/asb"
)

func (ui *UI) addConnection() {

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

			ui.PrintLog("Connection '" + newConnection.Name + "' added\n")
			ui.refreshConnections()
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		}).
		AddButton("Cancel", func() {
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		})

	ui.switchToForm("Add new connection")
}

func (ui *UI) removeSelectedConnection() {
	ui.advancedForm.clear()
	if ui.controller.GetSelectedConnection() == nil {
		ui.PrintLog("Please select connection to remove!\n")
		return
	}

	ui.advancedForm.setMessage(
		"Do you want to remove connection '" + ui.controller.GetSelectedConnection().Name + "'?",
	)
	ui.advancedForm.form.
		AddButton("Yes", func() {
			removed, err := ui.controller.RemoveSelectedConnection()
			if err != nil {
				ui.advancedForm.setMessage(err.Error())
				return
			}
			ui.PrintLog("Removed connection '" + removed + "'\n")
			ui.refreshConnections()
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		}).
		AddButton("No", func() {
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		})

	ui.switchToForm("Remove connection")
}

func (ui *UI) updateSelectedConnection() {
	ui.advancedForm.clear()
	newConnection := ui.controller.GetSelectedConnection()
	if newConnection == nil {
		ui.PrintLog("Please select connection to update!\n")
		return
	}
	ui.advancedForm.form.
		AddInputField("Connection name", newConnection.Name, 0, nil, nil).
		AddInputField("Service Bus namespace", newConnection.Namespace, 0, nil, nil).
		AddButton("Update", func() {
			ui.advancedForm.message.Clear()
			newConnection.Name = ui.advancedForm.form.GetFormItem(0).(*tview.InputField).GetText()
			newConnection.Namespace = ui.advancedForm.form.GetFormItem(1).(*tview.InputField).GetText()

			err := ui.controller.UpdateSelectedConnection(newConnection)
			if err != nil {
				ui.advancedForm.setMessage(err.Error())
				return
			}

			ui.PrintLog("Connection updated\n")
			ui.refreshConnections()
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		}).
		AddButton("Cancel", func() {
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		})

	ui.switchToForm("Update connection")
}
