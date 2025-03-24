// Package nop transports nothing.
package nop

import (
	"github.com/landru29/cnc-serial/internal/transport"
)

const defaultStatus = "<Idle|MPos:0.000,0.000,0.000|WPos:0.000,0.000,0.000|FS:0,0|Pn:P>"

var _ transport.TransportCloser = &Client{}

// Client is a serial client for sending commands.
type Client struct{}

// New creates the client.
func New() *Client {
	return &Client{}
}

// Send implements the Transport.Transporter interface.
func (c *Client) Send(_ ...string) error {
	return nil
}

// Close implements the io.Closer interface.
func (c Client) Close() error {
	return nil
}

// Read implements the io.Reader interface.
func (c Client) Read(data []byte) (int, error) {
	for idx, character := range defaultStatus + "\n" {
		data[idx] = byte(character)
	}

	return len(defaultStatus) + 1, nil
}
