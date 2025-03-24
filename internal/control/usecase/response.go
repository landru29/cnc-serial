package usecase

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/landru29/cnc-serial/internal/model"
)

func (c *Controller) bind(ctx context.Context) { //nolint: gocognit
	bufferline := ""

	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.transporterSetMutex.Lock()

			if c.transporter == nil {
				c.transporterSetMutex.Unlock()
				continue
			}

			buf := make([]byte, bufferSize)

			count, err := c.transporter.Read(buf)
			c.transporterSetMutex.Unlock()

			switch {
			case errors.Is(err, io.EOF):
				// Do nothing
			case err != nil:
				for _, display := range c.displayList {
					model.NewResponse(err.Error(), true).Encode(display)
				}
			default:
				bufferline += string(buf[:count])

				lineSplitter := strings.Split(bufferline, "\n")

				if len(lineSplitter) > 1 {
					for idx := 0; idx < len(lineSplitter)-1; idx++ {
						out := c.processResponse(lineSplitter[idx])
						if out == "" {
							continue
						}

						for _, display := range c.displayList {
							model.NewResponse(out, false).Encode(display)
						}
					}

					bufferline = lineSplitter[len(lineSplitter)-1]
				}
			}

			time.Sleep(delayBetweenSerialReads)
		}
	}
}

func (c *Controller) processResponse(resp string) string {
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
			_ = c.PushCommands(cmd)
		} else {
			c.status.CanRun = false
		}

		c.status.RemainingProgram = uint64(len(c.commandsToLaunch.commands) / 2)
	}

	_ = c.displayStatus()

	return ""
}
