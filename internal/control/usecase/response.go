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
			_ = model.NewResponse(err.Error(), true).Encode(display)
		}
	default:
		c.bufferline = append(c.bufferline, data...)

		lineSplitter := bytes.Split(c.bufferline, []byte("\n"))

		if len(lineSplitter) > 1 {
			for idx := 0; idx < len(lineSplitter)-1; idx++ {
				out := c.processLine(ctx, strings.TrimSpace(string(lineSplitter[idx])))
				if strings.TrimSpace(out) == "" {
					continue
				}

				for _, display := range c.displayList {
					_ = model.NewResponse(out, false).Encode(display)
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
		cmd, err := c.programmer.NextCommandToExecute()
		if err == nil {
			_ = c.PushCommands(ctx, true, cmd)
			_ = c.PushCommands(ctx, true, c.processer.CommandStatus())

			_ = c.displayProgram()
		}
	}

	_ = c.displayStatus()

	return ""
}

func (c *Controller) displayProgram() error {
	progModel := c.programmer.ToModel()
	for _, display := range c.displayList {
		if progModel != nil {
			if err := progModel.Encode(display); err != nil {
				return err
			}
		}
	}

	return nil
}
