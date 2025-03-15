// Package display manage the layout display.
package display

import (
	"io"

	"github.com/landru29/serial/internal/control"
	"github.com/rivo/tview"
)

// Screen is the main layout.
type Screen struct {
	xButton axisButtons
	yButton axisButtons
	zButton axisButtons

	display   *tview.Application
	userInput *tview.InputField
	logArea   *tview.TextView
	helpArea  *tview.TextView

	commander control.Commander
}

// New creates a screen.
func New(help func(string) string) *Screen {
	var output Screen

	output.buildView(help)

	return &output
}

// Start launches the tview application.
func (s Screen) Start() error {
	return s.display.EnableMouse(true).Run()
}

// Output is the main log area for displaying commands.
func (s Screen) Output() io.Writer {
	return s.logArea
}

// SetCommander sets the way to send commands.
func (s *Screen) SetCommander(commander control.Commander) {
	s.commander = commander
}
