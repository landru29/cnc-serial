// Package nop is the "no operation" implementation to send commands.
package nop

import (
	"fmt"
	"io"
	"strings"

	"github.com/landru29/serial/internal/gcode"
	"github.com/landru29/serial/internal/stack"
)

// Client is a serial client for sending commands.
type Client struct {
	coordinateRelative bool
	display            []io.Writer
	commandStack       *stack.Stack
}

// New creates the client.
func New(display ...io.Writer) *Client {
	return &Client{
		commandStack: &stack.Stack{},
		display:      display,
	}
}

// AddDisplay append display methods.
func (c *Client) AddDisplay(display ...io.Writer) {
	c.display = append(c.display, display...)
}

// Send implements the control.Commander interface.
func (c *Client) Send(texts ...string) error {
	for _, text := range texts {
		for _, display := range c.display {
			_, _ = fmt.Fprintf(display, " > %s\n", gcode.Colorize(text))
		}

		switch strings.ToUpper(strings.Split(text, " ")[0]) {
		case "G91":
			c.coordinateRelative = true
		case "G90":
			c.coordinateRelative = false
		}

		c.commandStack.Push(strings.ToUpper(text))
	}

	return nil
}

// IsRelative implements the control.Commander interface.
func (c *Client) IsRelative() bool {
	return c.coordinateRelative
}

// Close implements the io.Closer interface.
func (c Client) Close() error {
	return nil
}

// Bind implements the control.Commander interface.
func (c Client) Bind() {
}

// CommandStack implements the control.Commander interface.
func (c Client) CommandStack() *stack.Stack {
	return c.commandStack
}
