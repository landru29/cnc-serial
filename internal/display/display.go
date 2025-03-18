// Package display manage the layout display.
package display

import (
	"encoding/json"
	"fmt"

	"github.com/landru29/cnc-serial/internal/control"
	"github.com/landru29/cnc-serial/internal/gcode"
	"github.com/landru29/cnc-serial/internal/lang"
	"github.com/landru29/cnc-serial/internal/model"
	"github.com/landru29/cnc-serial/internal/stack"
	"github.com/rivo/tview"
)

// BaseScreen is the main layout.
type BaseScreen struct {
	display    *tview.Application
	userInput  *tview.InputField
	logArea    *tview.TextView
	helpArea   *tview.TextView
	statusArea *tview.TextView

	commander      control.Commander
	stackRetriever stack.Retriever
	currentLang    lang.Language
}

// New creates a screen.
func New(stackRetriever stack.Retriever, processer gcode.Processor) *Screen {
	output := Screen{
		BaseScreen: BaseScreen{
			stackRetriever: stackRetriever,
		},
	}

	output.buildView(processer)

	return &output
}

// Start launches the tview application.
func (s *Screen) Start() error {
	return s.display.EnableMouse(true).Run()
}

// Write implements the io.Writer interface.
func (s Screen) Write(p []byte) (n int, err error) {
	var status model.Status
	if err := json.Unmarshal(p, &status); err == nil {
		text := fmt.Sprintf("X: %03.1f\t\tY: %03.1f\t\tZ: %03.1f", status.XCoordinate, status.YCoordinate, status.ZCoordinate)
		s.statusArea.SetText(text)

		return len(text), nil
	}

	return s.logArea.Write(p)
}

// SetCommandSender sets the way to send commands.
func (s *Screen) SetCommandSender(commander control.Commander) {
	s.commander = commander
}

// SetLanguage sets the language.
func (s *Screen) SetLanguage(language lang.Language) {
	s.currentLang = language
}
