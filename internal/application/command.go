package application

import (
	"fmt"

	"github.com/landru29/serial/internal/gcode"
)

// SendCommand sends a command to the serial port.
func (c *Client) SendCommand(text string) {
	_, _ = fmt.Fprintf(c.logArea, " > %s\n", gcode.Colorize(text))

	if !c.dryRun() {
		if _, err := fmt.Fprintf(c.port, "%s\n", text); err != nil {
			_, _ = fmt.Fprintf(c.logArea, " - ERR %s\n", err.Error())
		}
	}

	c.lastCommand = append(c.lastCommand, text)
	c.cursor = len(c.lastCommand) - 1
}
