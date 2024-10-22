package ui

import (
	"github.com/rivo/tview"
)

type AdvancedForm struct {
	form    *tview.Form
	message *tview.TextView

	theme Theme
	flex  *tview.Flex
}

func newAdvancedForm(theme Theme) *AdvancedForm {

	form := tview.NewForm()
	form.SetBackgroundColor(theme.backgroundColor)
	form.SetButtonStyle(theme.style)
    form.SetFieldBackgroundColor(theme.foregroundColor).
        SetFieldTextColor(theme.backgroundColor)

	message := tview.NewTextView()
	message.SetTextStyle(theme.error).
        SetBackgroundColor(theme.backgroundColor)

	Flex := tview.NewFlex().SetDirection(tview.FlexRow)
	Flex.
		SetBorder(true).
		SetBackgroundColor(theme.backgroundColor)

	Flex.
		AddItem(message, 1, 0, false).
		AddItem(form, 0, 1, false)

	return &AdvancedForm{
		form:    form,
		message: message,

		theme: theme,
		flex:  Flex,
	}
}

func (advancedForm *AdvancedForm) clear() {
	advancedForm.form.Clear(true)
}

func (advancedForm *AdvancedForm) setMessage(msg string) {
	advancedForm.message.Clear()
	advancedForm.message.SetTextStyle(advancedForm.theme.error)

	advancedForm.message.SetText("\t" + msg)
}
