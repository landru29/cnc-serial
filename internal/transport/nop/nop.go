// Package nop transports nothing.
package nop

import (
	"io"

	"github.com/landru29/cnc-serial/internal/transport"
)

var _ transport.TransportCloser = &Client{}

// Client is a serial client for sending commands.
type Client struct{}

// New creates the client.
func New() *Client {
	return &Client{}
}

// Send implements the Transport.Transporter interface.
func (c *Client) Send(texts ...string) error {
	return nil
}

// Close implements the io.Closer interface.
func (c Client) Close() error {
	return nil
}

// Read implements the io.Reader interface.
func (c Client) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}
