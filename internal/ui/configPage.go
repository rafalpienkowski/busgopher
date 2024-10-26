package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/rafalpienkowski/busgopher/internal/controller"
)

type ConfigPage struct {
	theme      Theme
	controller *controller.Controller
	switchPage switchPageFunc
	closeApp   closeAppFunc

	flex   *tview.Flex
	config *tview.TextArea
	logs   *tview.TextView

	sending  *BoxButton
	validate *BoxButton
	save     *BoxButton
	close    *BoxButton

	inputs []tview.Primitive
}

func newConfigPage(theme Theme, closeApp closeAppFunc, switchPage switchPageFunc) *ConfigPage {

	flex := tview.NewFlex()

	config := tview.NewTextArea()
	logs := tview.NewTextView()
	sending := newBoxButton("To Sending Page")
	validate := newBoxButton("Validate")
	save := newBoxButton("Save")
	close := newBoxButton("Close")

	inputs := []tview.Primitive{
		config,
		sending,
		validate,
		save,
		close,
	}

	configPage := ConfigPage{
		theme:      theme,
		switchPage: switchPage,
		closeApp:   closeApp,

		flex:   flex,
		config: config,
		logs:   logs,

		sending:  sending,
		validate: validate,
		save:     save,
		close:    close,

		inputs: inputs,
	}

	configPage.configureAppearence()
	configPage.setLayout()

	return &configPage
}

func (configPage *ConfigPage) configureAppearence() {

	configPage.config.
		SetTitle(" Current config: ").
		SetBorder(true).
		SetBackgroundColor(configPage.theme.backgroundColor)

	configPage.logs.
		SetTitle(" Logs: ").
		SetBorder(true).
		SetBackgroundColor(configPage.theme.backgroundColor)

	configPage.flex.
		SetBorder(true).
		SetBackgroundColor(configPage.theme.backgroundColor).
		SetTitle("Configuration")
}

func (configPage *ConfigPage) setLayout() {

	actions := tview.NewFlex()
	actions.
		AddItem(tview.NewBox().SetBackgroundColor(configPage.theme.backgroundColor), 0, 1, false).
		AddItem(configPage.sending, configPage.sending.GetWidth(), 0, false).
		AddItem(configPage.validate, configPage.validate.GetWidth(), 0, false).
		AddItem(configPage.save, configPage.save.GetWidth(), 0, false).
		AddItem(configPage.close, configPage.close.GetWidth(), 0, false)

	main := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(configPage.config, 0, 3, false).
		AddItem(actions, 3, 0, false).
		AddItem(configPage.logs, 0, 1, false)

	configPage.flex.
		AddItem(main, 0, 3, false)
}

func (configPage *ConfigPage) loadData(controller *controller.Controller) {
	configPage.controller = controller
	configPage.setActions()
    json, _ := configPage.controller.GetConfigString()
    configPage.printConfig(json)
}

func (configPage *ConfigPage) setActions() {
	configPage.sending.SetSelectedFunc(func() {
        configPage.switchPage("sending")
	})

	configPage.validate.SetSelectedFunc(func() {
		err := configPage.controller.ValidateConfig()
		if err != nil {
			configPage.printError(err)
		}
	})

	configPage.save.SetSelectedFunc(func() {
		err := configPage.controller.SaveConfigJson(configPage.config.GetText())
		if err != nil {
			configPage.printError(err)
		}
	})
	configPage.close.SetSelectedFunc(func() {
		configPage.closeApp()
	})
}

func (configPage *ConfigPage) printConfig(currentConfig string) {
	configPage.config.SetText(currentConfig, false)
}

func (configPage *ConfigPage) printError(err error) {
	configPage.printLog(fmt.Sprintf(
		"[%v]: [Error] %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		err.Error(),
	))
}

func (configPage *ConfigPage) printLog(logMsg string) {
	fmt.Fprintf(configPage.logs, "%v", logMsg)

	getAvailableRows := func() int {
		_, _, _, height := configPage.logs.GetRect()

		return height - 2 // Minus border
	}

	configPage.logs.SetMaxLines(getAvailableRows())
}

func (configPage *ConfigPage) setAfterDrawFunc(focusedElement tview.Primitive) {

	configPage.config.SetBorderColor(tcell.ColorWhite)
	configPage.sending.SetBorderColor(tcell.ColorWhite)
	configPage.validate.SetBorderColor(tcell.ColorWhite)
	configPage.save.SetBorderColor(tcell.ColorWhite)
	configPage.close.SetBorderColor(tcell.ColorWhite)

	switch focusedElement {
	case configPage.config:
		configPage.config.SetBorderColor(tcell.ColorBlue)
	case configPage.sending:
		configPage.sending.SetBorderColor(tcell.ColorBlue)
	case configPage.validate:
		configPage.validate.SetBorderColor(tcell.ColorBlue)
	case configPage.save:
		configPage.save.SetBorderColor(tcell.ColorBlue)
	case configPage.close:
		configPage.close.SetBorderColor(tcell.ColorBlue)
	}
}
