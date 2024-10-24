package ui

import (
	"github.com/rivo/tview"

	"github.com/rafalpienkowski/busgopher/internal/asb"
)

func (ui *UI) addMessage() {
	ui.advancedForm.clear()
	ui.advancedForm.form.
		AddInputField("Name", "", 0, nil, nil).
		AddInputField("Correlation Id", "", 0, nil, nil).
		AddInputField("Message Id", "", 0, nil, nil).
		AddInputField("Reply to", "", 0, nil, nil).
		AddInputField("Subject", "", 0, nil, nil).
		//TODO: Add Custom properties
		AddTextArea("Body", "", 0, 30, 0, nil).
		AddButton("Add", func() {
			ui.advancedForm.message.Clear()
			newMessage := asb.Message{
				Name:          ui.advancedForm.form.GetFormItem(0).(*tview.InputField).GetText(),
				CorrelationID: ui.advancedForm.form.GetFormItem(1).(*tview.InputField).GetText(),
				MessageID:     ui.advancedForm.form.GetFormItem(2).(*tview.InputField).GetText(),
				ReplayTo:      ui.advancedForm.form.GetFormItem(3).(*tview.InputField).GetText(),
				Subject:       ui.advancedForm.form.GetFormItem(4).(*tview.InputField).GetText(),
				Body:          ui.advancedForm.form.GetFormItem(5).(*tview.TextArea).GetText(),
			}

			err := ui.controller.AddMessage(newMessage)
			if err != nil {
				ui.advancedForm.setMessage(err.Error())
				return
			}

			ui.PrintLog("Message '" + newMessage.Name + "' added\n")
			ui.refreshMessages()
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		}).
		AddButton("Cancel", func() {
			ui.pages.SwitchToPage("sending")
			ui.app.SetFocus(ui.connections)
		})

	ui.switchToForm("Add new message")
}

func (ui *UI) removeSelectedMessage() {
}

func (ui *UI) updateSelectedMessage() {
}
