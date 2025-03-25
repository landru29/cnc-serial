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

	commander      control.Commander
	stackRetriever stack.Retriever
	currentLang    lang.Language
	bufferData     string
	bufferMutex    sync.Mutex
}

// New creates a screen.
func New(ctx context.Context, stackRetriever stack.Retriever, processer gcode.Processor) *Screen {
	output := Screen{
		BaseScreen: BaseScreen{
			stackRetriever: stackRetriever,
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
}
