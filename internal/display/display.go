// Package display manage the layout display.
package display

import (
	"encoding/json"
	"fmt"
	"strings"

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
	bufferData     string
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
func (s *Screen) Write(data []byte) (int, error) {
	s.bufferData += string(data)

	splitter := strings.Split(s.bufferData, "\n")
	if len(splitter) < 2 { //nolint: mnd
		return len(data), nil
	}

	for _, line := range splitter {
		var status model.Status
		if err := json.Unmarshal([]byte(line), &status); err == nil {
			text := fmt.Sprintf(
				"X: %03.1f\t\tY: %03.1f\t\tZ: %03.1f",
				status.Machine.XCoordinate,
				status.Machine.YCoordinate,
				status.Machine.ZCoordinate,
			)
			s.statusArea.SetText(text)

			continue
		}

		_, _ = s.logArea.Write([]byte(line))
	}

	s.bufferData = splitter[len(splitter)-1]

	return len(data), nil
}

// SetCommandSender sets the way to send commands.
func (s *Screen) SetCommandSender(commander control.Commander) {
	s.commander = commander
}

// SetLanguage sets the language.
func (s *Screen) SetLanguage(language lang.Language) {
	s.currentLang = language
}
