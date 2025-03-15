package display

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (s *Screen) buildView(help func(string) string) {
	s.display = tview.NewApplication()

	s.logArea = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			s.display.Draw()
			s.logArea.ScrollToEnd()
		})

	s.helpArea = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			s.display.Draw()
		})

	s.userInput = tview.NewInputField().
		SetLabel("Enter command ").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite).
		SetDoneFunc(func(_ tcell.Key) {
			text := s.userInput.GetText()

			if strings.ToLower(text) == "exit" {
				s.display.Stop()

				return
			}

			_ = s.commander.Send(text)

			s.userInput.SetText("")
		})

	s.userInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() { //nolint: exhaustive
		case tcell.KeyUp:
			if cmd := s.commander.CommandStack().NavigateUp(); cmd != "" {
				s.userInput.SetText(cmd)
			}
		case tcell.KeyDown:
			if cmd := s.commander.CommandStack().NavigateDown(); cmd != "" {
				s.userInput.SetText(cmd)
			}

		default:
			s.commander.CommandStack().ResetCursor()
		}

		return event
	})

	s.userInput.SetChangedFunc(func(text string) {
		entry := strings.Split(text, " ")

		if description := help(strings.ToUpper(entry[0])); description != "" {
			s.helpArea.SetText(description)

			return
		}

		s.helpArea.SetText("")
	})

	s.logArea.SetBorder(true)
	s.helpArea.SetBorder(true)

	flex := tview.NewFlex().
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(s.logArea, 0, 8, false).  //nolint: mnd
				AddItem(s.helpArea, 0, 2, false). //nolint: mnd
				AddItem(s.userInput, 0, 1, true), 0, 4, true).
		AddItem(
			s.makeButtons(), 0, 1, false)

	s.display = s.display.SetRoot(flex, true)
}
