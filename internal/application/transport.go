package application

import (
	"github.com/landru29/cnc-serial/internal/control/usecase"
	"github.com/landru29/cnc-serial/internal/transport"
)

// SetTransport change the transport method to send commands.
func (c *Client) SetTransport(commander transport.TransportCloser) {
	c.screen.SetCommandSender(usecase.New(c.context, c.stack, commander, c.processer, c))
	c.transport = commander
}

// Close implements the io.Closer interface.
func (c *Client) Close() error {
	if c.transport != nil {
		err := c.transport.Close()
		if err != nil {
			return err
		}

		c.transport = nil
	}

	return nil
}
