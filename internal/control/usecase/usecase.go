// Package usecase is the control.Commander implementation.
package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/landru29/cnc-serial/internal/gcode"
	"github.com/landru29/cnc-serial/internal/model"
	"github.com/landru29/cnc-serial/internal/stack"
	"github.com/landru29/cnc-serial/internal/transport"
)

const (
	bufferSize = 200

	delayBetweenSerialReads = 500 * time.Millisecond

	delayBetweenStatusRequest = time.Second
)

// Controller is the control.Commander implementation.
type Controller struct {
	coordinateRelative bool
	display            []io.Writer
	stackPusher        stack.Pusher
	transporter        transport.Transporter
	processer          gcode.Processor
	mutex              sync.Mutex
	status             model.Status
}

// New creates the controller.
func New(
	ctx context.Context,
	stackPusher stack.Pusher,
	transporter transport.Transporter,
	processer gcode.Processor,
	display ...io.Writer,
) *Controller {
	output := &Controller{
		stackPusher: stackPusher,
		display:     display,
		transporter: transporter,
		processer:   processer,
	}

	go func() {
		output.bind(ctx)
	}()

	go func() {
		ticker := time.NewTicker(delayBetweenStatusRequest)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = output.PushCommands(processer.BuildStatusRequest())
			}
		}
	}()

	return output
}

// PushCommands implements the control.Commander interface.
func (c *Controller) PushCommands(commands ...string) error {
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
	}()

	for _, text := range commands {
		for _, display := range c.display {
			if text != c.processer.BuildStatusRequest() {
				_, _ = fmt.Fprintf(display, " > %s\n", c.processer.Colorize(text))
			}
		}

		if err := c.transporter.Send(text); err != nil {
			return err
		}

		switch strings.ToUpper(strings.Split(text, " ")[0]) {
		case "G91":
			c.coordinateRelative = true
		case "G90":
			c.coordinateRelative = false
		}

		if text != c.processer.BuildStatusRequest() {
			c.stackPusher.Push(strings.ToUpper(text))
		}
	}

	return nil
}

// IsRelative implements the control.Commander interface.
func (c *Controller) IsRelative() bool {
	return c.coordinateRelative
}

func (c *Controller) bind(ctx context.Context) { //nolint: gocognit
	bufferline := ""

	for {
		select {
		case <-ctx.Done():
			return
		default:
			buf := make([]byte, bufferSize)

			count, err := c.transporter.Read(buf)

			switch {
			case errors.Is(err, io.EOF):
				// Do nothing
			case err != nil:
				for _, display := range c.display {
					_, _ = fmt.Fprintf(display, " [#ff0000]ERR %s\n", err.Error())
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

						for _, display := range c.display {
							_, _ = fmt.Fprintf(display, " [#00ff00]%s\n", out)
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

	for _, display := range c.display {
		_ = json.NewEncoder(display).Encode(c.status) //nolint: errchkjson

		_, _ = display.Write([]byte("\n"))
	}

	return ""
}
