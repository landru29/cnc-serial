// Package serial is the serial implementation of the transport.Transporter interface.
package serial

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/landru29/cnc-serial/internal/transport"
	"go.bug.st/serial"
)

var _ transport.TransportCloser = &Client{}

const (
	bufferSize = 200

	delayBetweenSerialReads = 500 * time.Millisecond
)

// Client is a serial client for sending commands.
type Client struct {
	port         serial.Port
	handler      transport.ResponseHandler
	handlerMutex sync.Mutex
	stop         chan struct{}
	name         PortName
	bitRate      int
}

// New creates the client.
func New(ctx context.Context, name PortName, bitRate int) (*Client, error) {
	port, err := serial.Open(string(name), &serial.Mode{
		BaudRate: bitRate,
	})
	if err != nil {
		return nil, err
	}

	output := Client{
		port:    port,
		name:    name,
		bitRate: bitRate,
	}

	go func() {
		output.bind(ctx)
	}()

	return &output, nil
}

// ConnectionStatus implements the Transport.Transporter interface.
func (c *Client) ConnectionStatus() string {
	return fmt.Sprintf("%s@%d", c.name, c.bitRate)
}

// Send implements the Transport.Transporter interface.
func (c *Client) Send(_ context.Context, texts ...string) error {
	for _, text := range texts {
		if _, err := fmt.Fprintf(c.port, "%s\n", text); err != nil {
			return err
		}
	}

	return nil
}

// Close implements the io.Closer interface.
func (c *Client) Close() error {
	return c.port.Close()
}

// SetResponseHandler implements the Transport.Transporter interface.
func (c *Client) SetResponseHandler(handler transport.ResponseHandler) {
	c.handlerMutex.Lock()
	c.handler = handler
	c.handlerMutex.Unlock()
}

func (c *Client) bind(ctx context.Context) {
	defer close(c.stop)

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.stop:
			return
		default:
			c.handlerMutex.Lock()
			if c.handler != nil {
				data := make([]byte, bufferSize)
				count, err := c.port.Read(data)
				c.handler(ctx, data[:count], err)
			}
			c.handlerMutex.Unlock()

			time.Sleep(delayBetweenSerialReads)
		}
	}
}
