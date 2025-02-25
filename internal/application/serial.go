package application

import (
	"errors"
	"fmt"
	"io"
	"time"

	"go.bug.st/serial"
)

const (
	bufferSize = 200

	delayBetweenSerialReads = 500 * time.Millisecond
)

// OpenPort opens the serial port.
func (c *Client) OpenPort(name string, bitRate int) error {
	port, err := serial.Open(name, &serial.Mode{
		BaudRate: bitRate,
	})
	if err != nil {
		return err
	}

	c.port = port

	return nil
}

// Close implements the io.Closer interface.
func (c *Client) Close() error {
	if c.port != nil {
		err := c.port.Close()
		if err != nil {
			return err
		}

		c.port = nil
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

// Bind reads data from the serial port when available.
func (c Client) Bind() {
	for {
		buf := make([]byte, bufferSize)

		count, err := c.port.Read(buf)

		switch {
		case errors.Is(err, io.EOF):
			// Do nothing
		case err != nil:
			_, _ = fmt.Fprintf(c.logArea, " [#ff0000]ERR %s\n", err.Error())
		default:
			_, _ = fmt.Fprintf(c.logArea, " [#00ff00]%s", string(buf[:count]))
		}

		time.Sleep(delayBetweenSerialReads)
	}
}
