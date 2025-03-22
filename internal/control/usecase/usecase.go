// Package usecase is the control.Commander implementation.
package usecase

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/landru29/cnc-serial/internal/control"
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

var _ control.Commander = &Controller{}

// Controller is the control.Commander implementation.
type Controller struct {
	displayList         []io.Writer
	stackPusher         stack.Pusher
	transporter         transport.Transporter
	transporterSetMutex sync.Mutex
	processer           gcode.Processor
	pushMutex           sync.Mutex
	status              model.Status
	programmer          gcode.Programmer
	programmerSetMutex  sync.Mutex
}

// New creates the controller.
func New(
	ctx context.Context,
	stackPusher stack.Pusher,
	processer gcode.Processor,
	displayList ...io.Writer,
) *Controller {
	output := &Controller{
		stackPusher: stackPusher,
		displayList: displayList,
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
				_ = output.PushCommands(processer.CommandStatus())
			}
		}
	}()

	return output
}

func (c *Controller) SetTransporter(transporter transport.Transporter) {
	c.transporterSetMutex.Lock()
	defer func() {
		c.transporterSetMutex.Unlock()
	}()

	c.transporter = transporter
}

func (c *Controller) SetProgrammer(programmer gcode.Programmer) {
	c.programmerSetMutex.Lock()
	defer func() {
		c.programmerSetMutex.Unlock()
	}()

	c.programmer = programmer
}

// PushCommands implements the control.Commander interface.
func (c *Controller) PushCommands(commands ...string) error {
	c.pushMutex.Lock()
	c.transporterSetMutex.Lock()
	defer func() {
		c.pushMutex.Unlock()
		c.transporterSetMutex.Unlock()
	}()

	if c.transporter == nil {
		return errors.New("missing transporter")
	}

	for _, text := range commands {
		if strings.TrimSpace(text) == "" {
			if err := c.stepProgram(); err != nil {
				return err
			}

			continue
		}

		for _, display := range c.displayList {
			if text != c.processer.CommandStatus() {
				model.NewRequest(c.processer.Colorize(text)).Encode(display)
			}
		}

		if err := c.transporter.Send(text); err != nil {
			return err
		}

		switch strings.ToUpper(strings.Split(text, " ")[0]) {
		case c.processer.CommandRelativeCoordinate():
			c.status.RelativeCoordinates = true

			_ = c.displayStatus()

		case c.processer.CommandAbsoluteCoordinate():
			c.status.RelativeCoordinates = false

			_ = c.displayStatus()

		}

		if text != c.processer.CommandStatus() {
			c.stackPusher.Push(strings.ToUpper(text))
		}
	}

	return nil
}

func (c *Controller) stepProgram() error {
	c.programmerSetMutex.Lock()
	defer c.programmerSetMutex.Unlock()

	if c.programmer == nil {
		return nil
	}

	currentCommand := c.programmer.CurrentCommand()

	if _, err := c.programmer.ReadNextInstruction(); err != nil {
		return err
	}

	for _, display := range c.displayList {
		model.NewRequest(c.processer.Colorize(currentCommand)).Encode(display)

		model := c.programmer.ToModel()
		if model != nil {
			if err := model.Encode(display); err != nil {
				return err
			}
		}
	}

	if err := c.transporter.Send(currentCommand); err != nil {
		return err
	}

	return nil
}

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

	_ = c.displayStatus()

	return ""
}

func (c *Controller) displayStatus() error {
	for _, display := range c.displayList {
		if err := c.status.Encode(display); err != nil {
			return err
		}
	}

	return nil
}

// MoveRelative implements the control.Commander interface.
func (c *Controller) MoveRelative(offset float64, axisName string) error {
	commands := []string{
		c.processer.CommandRelativeCoordinate(),
		c.processer.MoveAxis(offset, axisName),
		c.processer.CommandAbsoluteCoordinate(),
	}

	if c.status.RelativeCoordinates {
		commands = commands[1:2]
	}

	return c.PushCommands(commands...)
}
