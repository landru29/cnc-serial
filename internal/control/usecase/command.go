package usecase

import (
	"context"
	"strconv"
	"strings"

	"github.com/landru29/cnc-serial/internal/errors"
	"github.com/landru29/cnc-serial/internal/model"
)

// PushCommands implements the control.Commander interface.
func (c *Controller) PushCommands(ctx context.Context, fromProgram bool, commands ...string) error { //nolint: funlen,gocognit,cyclop
	c.pushMutex.Lock()
	c.transporterSetMutex.Lock()
	defer func() {
		c.pushMutex.Unlock()
		c.transporterSetMutex.Unlock()
	}()

	if c.transporter == nil {
		return errors.ErrMissingTransporter
	}

	for _, text := range commands {
		command := strings.TrimSpace(text)

		if command == "" {
			if fromProgram {
				continue
			}

			c.status.CanRun = true

			c.programmer.SetLinesToExecute(1)

			continue
		}

		if command[0] == 's' || command[0] == 'S' {
			c.status.CanRun = false

			_ = c.displayStatus()

			continue
		}

		if command[0] == 'r' || command[0] == 'R' {
			c.status.CanRun = false

			_ = c.programmer.Reset()

			_ = c.displayProgram()

			continue
		}

		if command[0] == 'p' || command[0] == 'P' { //nolint: nestif
			if len(command) == 1 {
				c.status.CanRun = true

				_ = c.displayStatus()

				continue
			}

			cmdCount := int64(1)

			if len(command) > 1 {
				programMatcher := c.regexpProcessProgram.FindAllStringSubmatch(strings.TrimSpace(text), -1)
				if len(programMatcher) > 0 && len(programMatcher[0]) == 2 {
					if count, err := strconv.ParseInt(programMatcher[0][1], 10, 64); err == nil {
						cmdCount = count
					}
				}
			}

			c.programmer.SetLinesToExecute(cmdCount)

			c.status.CanRun = true

			continue
		}

		// Fix typing errors
		text = strings.ReplaceAll(strings.ToUpper(text), "O", "0")

		// Display.
		for _, display := range c.displayList {
			if text != c.processer.CommandStatus() {
				_ = model.NewRequest(c.processer.Colorize(text)).Encode(display)
			}
		}

		// Send command.
		if err := c.transporter.Send(ctx, text); err != nil {
			return err
		}

		// decode command.
		switch strings.ToUpper(strings.Split(text, " ")[0]) {
		case c.processer.CommandRelativeCoordinate():
			c.status.RelativeCoordinates = true

			_ = c.displayStatus()

		case c.processer.CommandAbsoluteCoordinate():
			c.status.RelativeCoordinates = false

			_ = c.displayStatus()
		}

		// Remember command.
		if !fromProgram && text != c.processer.CommandStatus() {
			c.stackPusher.Push(strings.ToUpper(text))
		}
	}

	return nil
}
