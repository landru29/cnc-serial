// Package nop transports nothing.
package nop

import (
	"context"
	"sync"
	"time"

	"github.com/landru29/cnc-serial/internal/transport"
)

var _ transport.TransportCloser = &Client{}

const (
	defaultStatus          = "<Idle|MPos:0.000,0.000,0.000|WPos:0.000,0.000,0.000|FS:0,0|Pn:P>"
	delayBetweenStatusSend = 500 * time.Millisecond
)

// Client is a serial client for sending commands.
type Client struct {
	handler      transport.ResponseHandler
	handlerMutex sync.Mutex
	stop         chan struct{}
}

// New creates the client.
func New(ctx context.Context) *Client {
	output := &Client{
		stop: make(chan struct{}, 1),
	}

	go func() {
		output.bind(ctx)
	}()

	return output
}

// ConnectionStatus implements the Transport.Transporter interface.
func (c *Client) ConnectionStatus() string {
	return "mock"
}

// Send implements the Transport.Transporter interface.
func (c *Client) Send(_ context.Context, _ ...string) error {
	return nil
}

// Close implements the io.Closer interface.
func (c *Client) Close() error {
	c.stop <- struct{}{}
	return nil
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
				c.handler(ctx, []byte(defaultStatus+"\n"), nil)
			}
			c.handlerMutex.Unlock()

			time.Sleep(delayBetweenStatusSend)
		}
	}
}
