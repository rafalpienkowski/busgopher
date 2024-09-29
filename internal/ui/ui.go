package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/rafalpienkowski/busgopher/internal/controller"
)

// UI implement terminal user interface features.
type UI struct {
	controller *controller.Controller

	// View components
	App          *tview.Application
	Flex         *tview.Flex
	Connections  *tview.List
	Messages     *tview.List
	Content      *tview.TextView
	Logs         *tview.TextView

	Inputs []tview.Primitive
}

func NewUI(controller *controller.Controller) *UI {
	ui := UI{}

	ui.controller = controller

	// Create UI elements
	ui.App = tview.NewApplication()
	ui.Connections = tview.NewList().
		ShowSecondaryText(false).
		SetWrapAround(true).
		SetHighlightFullLine(true)
	ui.Messages = tview.NewList().
		ShowSecondaryText(false).
		SetWrapAround(true).
		SetHighlightFullLine(true)
	ui.Content = tview.NewTextView()
	ui.Logs = tview.NewTextView()

	ui.Inputs = []tview.Primitive{
		ui.Connections,
		ui.Messages,
		ui.Content,
	}

	// Configure appearence
	ui.Connections.SetTitle(" Connections: ").SetBorder(true)
	ui.Messages.SetTitle(" Messages: ").SetBorder(true)
	ui.Content.SetTitle(" Content: ").SetBorder(true)
	ui.Logs.SetTitle(" Logs: ").SetBorder(true)

	// Configure handlers

	// Set layouts
	left := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.Connections, 0, 1, true).
		AddItem(ui.Messages, 0, 1, false)

	right := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.Content, 0, 3, false).
		AddItem(ui.Logs, 0, 1, false)

	ui.Flex = tview.NewFlex().
		AddItem(left, 0, 1, false).
		AddItem(right, 0, 3, false)

	ui.App.SetAfterDrawFunc(ui.setAfterDrawFunc)

	ui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			ui.cycleFocus(false)
		} else if event.Key() == tcell.KeyBacktab {
			ui.cycleFocus(true)
		}
		return event
	})

	return &ui
}

func (ui *UI) LoadData() {
	for _, conn := range ui.controller.Connections {
		ui.Connections.AddItem(conn.Name, conn.Namespace, 0, func(){
            ui.controller.SelectConnection(conn)
            ui.printLog("Selected connection: " + conn.Name + " (" + conn.Namespace + ")")
        })
	}

	for _, msg := range ui.controller.Messages {
		ui.Messages.AddItem(msg.Name, msg.Subject, 0, func() {
            ui.controller.SelectMessage(msg)
			ui.printContent(msg.Print())
            ui.printLog("Selected message: " + msg.Name)
		})
	}
}

func (ui *UI) Start() error {
	return ui.App.SetRoot(ui.Flex, true).SetFocus(ui.Connections).EnableMouse(true).Run()
}

func (ui *UI) printLog(logMsg string) {
	fmt.Fprintf(ui.Logs, "[%v]: %v\n", time.Now().Format("2006-01-02 15:04:05"), logMsg)
}

func (ui *UI) printContent(content string) {
	ui.Content.Clear()
	fmt.Fprintf(ui.Content, "%v", content)
}

// Changes focus on TAB pressed
func (ui *UI) cycleFocus(
	reverse bool,
) {
	for i, el := range ui.Inputs {
		if !el.HasFocus() {
			continue
		}

		if reverse {
			i = i - 1
			if i < 0 {
				i = len(ui.Inputs) - 1
			}
		} else {
			i = i + 1
			i = i % len(ui.Inputs)
		}
		ui.App.SetFocus(ui.Inputs[i])
		return
	}
}

func (ui *UI) queueUpdateDraw(f func()) {
	go func() {
		ui.App.QueueUpdateDraw(f)
	}()
}

func (ui *UI) setAfterDrawFunc(screen tcell.Screen) {
	ui.queueUpdateDraw(func() {
		p := ui.App.GetFocus()

		ui.Connections.SetBorderColor(tcell.ColorWhite)
		ui.Messages.SetBorderColor(tcell.ColorWhite)
		ui.Content.SetBorderColor(tcell.ColorWhite)
		ui.Logs.SetBorderColor(tcell.ColorWhite)

		switch p {
		case ui.Connections:
			ui.Connections.SetBorderColor(tcell.ColorBlue)
		case ui.Messages:
			ui.Messages.SetBorderColor(tcell.ColorBlue)
		case ui.Content:
			ui.Content.SetBorderColor(tcell.ColorBlue)
		case ui.Logs:
			ui.Logs.SetBorderColor(tcell.ColorBlue)
		}
	})
}
