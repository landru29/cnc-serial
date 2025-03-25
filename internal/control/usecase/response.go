package usecase

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"

	"github.com/landru29/cnc-serial/internal/model"
)

// ProcessResponse implements the control.Commander interface.
func (c *Controller) ProcessResponse(ctx context.Context, data []byte, err error) {
	switch {
	case errors.Is(err, io.EOF):
		// Do nothing
	case err != nil:
		for _, display := range c.displayList {
			model.NewResponse(err.Error(), true).Encode(display)
		}
	default:
		c.bufferline = append(c.bufferline, data...)

		lineSplitter := bytes.Split(c.bufferline, []byte("\n"))

		if len(lineSplitter) > 1 {
			for idx := 0; idx < len(lineSplitter)-1; idx++ {
				out := c.processLine(ctx, string(lineSplitter[idx]))
				if out == "" {
					continue
				}

				for _, display := range c.displayList {
					model.NewResponse(out, false).Encode(display)
				}
			}

			c.bufferline = lineSplitter[len(lineSplitter)-1]
		}
	}
}

func (c *Controller) processLine(ctx context.Context, resp string) string {
	if strings.TrimSpace(resp) == "ok" {
		return ""
	}

	status, err := c.processer.UnmarshalStatus(resp)
	if err != nil {
		return resp + "\n"
	}

	c.status.Merge(*status)

	// launch next program commands if available.
	if strings.ToUpper(string(c.status.State)) == "IDLE" && c.status.CanRun {
		cmd, found := c.commandsToLaunch.next()
		if found {
			_ = c.PushCommands(ctx, cmd)
		} else {
			c.status.CanRun = false
		}

		c.status.RemainingProgram = uint64(len(c.commandsToLaunch.commands) / 2)
	}

	_ = c.displayStatus()

	return ""
}
