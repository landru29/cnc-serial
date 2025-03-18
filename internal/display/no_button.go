//go:build !withbutton

package display

import (
	"github.com/rivo/tview"
)

// Screen is the main layout.
type Screen struct {
	BaseScreen
}

func (s Screen) layout() *tview.Flex {
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(s.statusArea, 1, 0, false).
		AddItem(s.logArea, 0, 8, false).  //nolint: mnd
		AddItem(s.helpArea, 0, 2, false). //nolint: mnd
		AddItem(s.userInput, 0, 1, true)
}
