package ui

import "github.com/gdamore/tcell/v2"

type Theme struct {
	backgroundColor tcell.Color
	foregroundColor tcell.Color
	style           tcell.Style
}

func Dark() Theme {
	return Theme{
		backgroundColor: tcell.ColorGray,
		foregroundColor: tcell.ColorWhite,
		style: tcell.StyleDefault.Background(tcell.ColorGray).
			Foreground(tcell.ColorWhite),
	}
}
