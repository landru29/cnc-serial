// Package serial is the serial implementation of the transport.Transporter interface.
package serial

import (
	"fmt"

	"github.com/landru29/cnc-serial/internal/transport"
	"go.bug.st/serial"
)

var _ transport.TransportCloser = &Client{}

// Client is a serial client for sending commands.
type Client struct {
	port serial.Port
}

// New creates the client.
func New(name string, bitRate int) (*Client, error) {
	port, err := serial.Open(name, &serial.Mode{
		BaudRate: bitRate,
	})
	if err != nil {
		return nil, err
	}

	output := Client{
		port: port,
	}

	return &output, nil
}

// Send implements the Transport.Transporter interface.
func (c *Client) Send(texts ...string) error {
	for _, text := range texts {
		if _, err := fmt.Fprintf(c.port, "%s\n", text); err != nil {
			return err
		}
	}

	return nil
}

// Close implements the io.Closer interface.
func (c Client) Close() error {
	return c.port.Close()
}

// Read implements the io.Reader interface.
func (c Client) Read(p []byte) (int, error) {
	return c.port.Read(p)
}
