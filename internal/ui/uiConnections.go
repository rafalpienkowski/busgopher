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
            
            ui.PrintLog("Connection '" + newConnection.Name + "' added")
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

func (ui *UI) removeConnection() {
	ui.advancedForm.clear()
    if len(ui.controller.GetSelectedConnectionName()) == 0 {
        ui.PrintLog("Please select connection to remove!")
        return 
    }

    ui.advancedForm.setMessage("Do you want to remove connection '" + ui.controller.GetSelectedConnectionName() + "'?")
	ui.advancedForm.form.
		AddButton("Yes", func() {
			removed, err := ui.controller.RemoveSelectedConnection()
			if err != nil {
				ui.advancedForm.setMessage(err.Error())
				return
			}
            ui.PrintLog("Removed connection '" + removed + "'")
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
    if len("ui.controller.SelectedConnectionName") == 0 {
        ui.PrintLog("Please select connection to update")
        return
    }

	ui.switchToForm("Update connection")
}
