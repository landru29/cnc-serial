package display

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/landru29/cnc-serial/internal/gcode"
	"github.com/landru29/cnc-serial/internal/gpm"
	"github.com/rivo/tview"
)

const (
	enterCommandLabel = "Enter command"
)

func (s *Screen) buildView(ctx context.Context, processer gcode.Processor) {
	s.processer = processer

	s.display = tview.NewApplication().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			s.display.Stop()

			return event
		}

		if event.Modifiers()&tcell.ModCtrl != 0 {
			_ = s.commander.PushCommands(ctx, false, fmt.Sprintf("0x%0x\n", event.Key()))
		}

		return event
	})

	screen, err := gpm.NewScreen()

	if err == nil {
		s.display.SetScreen(screen)
	}

	s.logArea = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			s.display.Draw()
			s.logArea.ScrollToEnd()
		})

	s.progArea = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			s.display.Draw()
			s.logArea.ScrollToEnd()
		})

	s.statusArea = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(false).
		SetChangedFunc(func() {
			s.display.Draw()
		})

	s.helpArea = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			s.display.Draw()
		})

	s.userInput = tview.NewInputField().
		SetLabel(enterCommandLabel + " ").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite).
		SetDoneFunc(func(_ tcell.Key) {
			text := s.userInput.GetText()

			if strings.ToLower(text) == "exit" {
				s.display.Stop()

				return
			}

			_ = s.commander.PushCommands(ctx, false, text)

			s.userInput.SetText("")
		})

	s.userInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() { //nolint: exhaustive
		case tcell.KeyUp:
			if cmd := s.stackRetriever.NavigateUp(); cmd != "" {
				s.userInput.SetText(cmd)
			}
		case tcell.KeyDown:
			if cmd := s.stackRetriever.NavigateDown(); cmd != "" {
				s.userInput.SetText(cmd)
			}

		default:
			s.stackRetriever.ResetCursor()
		}

		return event
	})

	s.userInput.SetChangedFunc(func(text string) {
		if description := processer.CodeDescription(s.currentLang, strings.Split(text, " ")[0]); description != "" {
			s.helpArea.SetText(description)

			return
		}

		if description := processer.CodeDescription(s.currentLang, gcode.DefaultHelperCode); description != "" {
			s.helpArea.SetText(description)

			return
		}
	})

	s.logArea.SetBorder(true)
	s.helpArea.SetBorder(true)
	s.progArea.SetBorder(true)

	s.display = s.display.SetRoot(s.layout(ctx), true)
}
