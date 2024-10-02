package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type BoxButton struct {
	*tview.Box
	label   string
	focused bool
}

func (bb *BoxButton) NewBoxButton(label string) *BoxButton {
	box := tview.NewBox().
		SetBorder(true).
		SetTitleAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorGray).
		SetBorderColor(tcell.ColorWhite).
		SetBorderAttributes(tcell.AttrBold).
		SetTitleColor(tcell.ColorWhite)

	return &BoxButton{
		Box:   box,
		label: label,
	}
}

func (b *BoxButton) Draw(screen tcell.Screen) {
	b.Box.DrawForSubclass(screen, b)

	x, y, width, height := b.GetInnerRect()

	textWidth := len(b.label)
	textX := x + (width-textWidth)/2
	textY := y + height/2

	// Draw the label centered in the box
	for i, char := range b.label {
		screen.SetContent(
			textX+i,
			textY,
			char,
			nil,
			tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGray),
		)
	}
}

func (b *BoxButton) GetWidth() int {
    return len(b.label) + 4
}

// InputHandler returns the handler for input events (e.g., simulating button clicks).
func (b *BoxButton) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(tview.Primitive)) {
		if event.Key() == tcell.KeyEnter {
			// Simulate a "click" by changing the label
			b.label = "Clicked!"
		}
	}
}

func (b *BoxButton) Focus(delegate func(tview.Primitive)) {
	b.focused = true
}

func (b *BoxButton) Blur() {
	b.focused = false
}

func (b *BoxButton) HasFocus() bool {
	return b.focused
}
