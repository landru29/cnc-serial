// Package application is the main application.
package application

import (
	"context"
	"io"
	"strings"

	"github.com/landru29/cnc-serial/internal/control"
	"github.com/landru29/cnc-serial/internal/control/usecase"
	"github.com/landru29/cnc-serial/internal/display"
	"github.com/landru29/cnc-serial/internal/gcode"
	"github.com/landru29/cnc-serial/internal/lang"
	"github.com/landru29/cnc-serial/internal/model"
	"github.com/landru29/cnc-serial/internal/stack"
)

// Client is the main application structure.
type Client struct {
	context   context.Context //nolint: containedctx
	transport io.Closer
	stack     stack.Stacker
	screen    *display.Screen
	processer gcode.Processor
	program   gcode.Programmer
	commander control.Commander
}

// NewClient initializes a new application client.
func NewClient(
	ctx context.Context,
	stacker stack.Stacker,
	processer gcode.Processor,
	program gcode.Programmer,
	navigationInc float64,
) (*Client, error) {
	screen := display.New(ctx, stacker, processer, navigationInc)

	output := &Client{
		stack:     stacker,
		processer: processer,
		screen:    screen,
		context:   ctx,
		program:   program,
	}

	output.commander = usecase.New(ctx, stacker, processer, output)

	if program != nil {
		if err := program.ReadNextInstruction(); err != nil {
			return nil, err
		}

		output.commander.SetProgrammer(program)

		model := program.ToModel()
		if model != nil {
			if err := model.Encode(screen); err != nil {
				return nil, err
			}
		}
	}

	if err := (model.Status{
		Machine:    &model.Coordinates{},
		ToolOffset: &model.Coordinates{},
		State:      model.State("Waiting"),
		CanRun:     true,
	}.Encode(screen)); err != nil {
		return nil, err
	}

	return output, nil
}

// SetLanguage sets the language.
func (c *Client) SetLanguage(language lang.Language) {
	c.screen.SetLanguage(lang.Language(strings.ToLower(string(language))))
}

// AvailableLanguages lists all available languages.
func (c Client) AvailableLanguages() []lang.Language {
	return c.processer.AvailableLanguages()
}

// Start launches the tview application.
func (c Client) Start() error {
	return c.screen.Start()
}

// Write implements the io.Writer interface.
func (c Client) Write(p []byte) (int, error) {
	return c.screen.Write(p)
}
