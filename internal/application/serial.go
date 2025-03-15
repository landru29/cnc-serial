package application

import (
	controlserial "github.com/landru29/serial/internal/control/serial"
	"go.bug.st/serial"
)

// OpenPort opens the serial port.
func (c *Client) OpenPort(name string, bitRate int) error {
	port, err := serial.Open(name, &serial.Mode{
		BaudRate: bitRate,
	})
	if err != nil {
		return err
	}

	c.commander = controlserial.New(port, c.screen.Output())

	c.screen.SetCommander(c.commander)

	return nil
}

// Close implements the io.Closer interface.
func (c *Client) Close() error {
	if c.commander != nil {
		err := c.commander.Close()
		if err != nil {
			return err
		}

		c.commander = nil
	}

	return nil
}

// DefaultPort lists all available ports and chooses one as the default port to open.
func (c Client) DefaultPort() string {
	ports, err := serial.GetPortsList()
	if err == nil && len(ports) > 0 {
		return ports[0]
	}

	return ""
}
