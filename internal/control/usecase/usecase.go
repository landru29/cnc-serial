// Package usecase is the control.Commander implementation.
package usecase

import (
	"context"
	"io"
	"regexp"
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
	commandsToLaunch    *commandSet

	regexpProcessProgram *regexp.Regexp
}

// New creates the controller.
func New(
	ctx context.Context,
	stackPusher stack.Pusher,
	processer gcode.Processor,
	displayList ...io.Writer,
) *Controller {

	output := &Controller{
		stackPusher:          stackPusher,
		displayList:          displayList,
		processer:            processer,
		commandsToLaunch:     &commandSet{},
		regexpProcessProgram: regexp.MustCompile(`(?i)p([-+\d]+)`),
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
