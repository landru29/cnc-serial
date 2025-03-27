package usecase

import "io"

type commandSet struct {
	commands []string
}

func (c *commandSet) next() (string, bool) {
	if len(c.commands) == 0 {
		return "", false
	}

	output := c.commands[0]
	c.commands = c.commands[1:]

	return output, true
}

func (c *commandSet) push(statusCommand string, commands ...string) {
	for _, command := range commands {
		c.commands = append(c.commands, statusCommand, command)
	}
}

func (c *Controller) stepProgram(count int64) error {
	c.programmerSetMutex.Lock()
	defer c.programmerSetMutex.Unlock()

	if c.programmer == nil {
		return nil
	}

	for range count {
		currentCommand := c.programmer.CurrentCommand()

		if currentCommand == "" {
			return io.EOF
		}

		if _, err := c.programmer.ReadNextInstruction(); err != nil {
			return err
		}

		if err := c.PushProgramCommands(currentCommand); err != nil {
			return err
		}
	}

	progModel := c.programmer.ToModel()
	for _, display := range c.displayList {
		if progModel != nil {
			if err := progModel.Encode(display); err != nil {
				return err
			}
		}
	}

	return c.displayStatus()
}

// PushProgramCommands implements the control.Commander interface.
func (c *Controller) PushProgramCommands(commands ...string) error {
	c.status.RemainingProgram = int64(len(c.commandsToLaunch.commands) / 2) //nolint: mnd
	c.commandsToLaunch.push(c.processer.CommandStatus(), commands...)

	if err := c.displayStatus(); err != nil {
		return err
	}

	return nil
}
