// Package display manage the layout display.
package display

import (
	"context"
	"sync"

	"github.com/landru29/cnc-serial/internal/control"
	"github.com/landru29/cnc-serial/internal/gcode"
	"github.com/landru29/cnc-serial/internal/lang"
	"github.com/landru29/cnc-serial/internal/stack"
	"github.com/rivo/tview"
)

// BaseScreen is the main layout.
type BaseScreen struct {
	display    *tview.Application
	userInput  *tview.InputField
	logArea    *tview.TextView
	progArea   *tview.TextView
	helpArea   *tview.TextView
	statusArea *tview.TextView

	processer gcode.Processor

	commander      control.Commander
	stackRetriever stack.Retriever
	currentLang    lang.Language
	bufferData     string
	bufferMutex    sync.Mutex
	navigationInc  float64
}

// New creates a screen.
func New(ctx context.Context, stackRetriever stack.Retriever, processer gcode.Processor, navigationInc float64) *Screen {
	output := Screen{
		BaseScreen: BaseScreen{
			stackRetriever: stackRetriever,
			navigationInc:  navigationInc,
		},
	}

	output.buildView(ctx, processer)

	return &output
}

// Start launches the tview application.
func (s *Screen) Start() error {
	return s.display.EnableMouse(true).Run()
}

// SetCommandSender sets the way to send commands.
func (s *Screen) SetCommandSender(commander control.Commander) {
	s.commander = commander
}

// SetLanguage sets the language.
func (s *Screen) SetLanguage(language lang.Language) {
	s.currentLang = language

	if s.processer == nil {
		return
	}

	if description := s.processer.CodeDescription(s.currentLang, gcode.DefaultHelperCode); description != "" {
		s.helpArea.SetText(description)

		return
	}
}
