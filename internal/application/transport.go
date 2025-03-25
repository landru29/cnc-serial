package application

import (
	"github.com/landru29/cnc-serial/internal/transport"
)

// SetTransport change the transport method to send commands.
func (c *Client) SetTransport(transporter transport.TransportCloser) {
	c.commander.SetTransporter(transporter)
	c.screen.SetCommandSender(c.commander)
	transporter.SetResponseHandler(c.commander.ProcessResponse)
	c.transport = transporter
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
