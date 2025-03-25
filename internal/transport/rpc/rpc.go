package rpc

import (
	"context"
	"sync"
	"time"

	"github.com/landru29/cnc-serial/internal/transport"
	rpcmodel "github.com/landru29/cnc-serial/internal/transport/rpc/model"
	"google.golang.org/grpc"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative model/message.proto

var _ transport.TransportCloser = &Client{}

const (
	delayBetweenStatusSend = 500 * time.Millisecond
)

// Client is a serial client for sending commands.
type Client struct {
	handler      transport.ResponseHandler
	handlerMutex sync.Mutex
	stop         chan struct{}
	conn         *grpc.ClientConn
	client       rpcmodel.CommandSenderClient
}

// New creates the client.
func New(ctx context.Context, serverAddr string, opts ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		return nil, err
	}

	output := &Client{
		stop:   make(chan struct{}, 1),
		conn:   conn,
		client: rpcmodel.NewCommandSenderClient(conn),
	}

	go func() {
		output.bind(ctx)
	}()

	return output, nil
}

// Send implements the Transport.Transporter interface.
func (c *Client) Send(ctx context.Context, commands ...string) error {
	for _, cmd := range commands {
		if _, err := c.client.SendCommand(ctx, &rpcmodel.Command{
			Data: cmd,
		}); err != nil {
			return err
		}
	}

	return nil
}

// Close implements the io.Closer interface.
func (c *Client) Close() error {
	c.stop <- struct{}{}

	return c.conn.Close()
}

// SetResponseHandler implements the Transport.Transporter interface.
func (c *Client) SetResponseHandler(handler transport.ResponseHandler) {
	c.handler = handler
}

func (c *Client) reply(ctx context.Context, data string) {
	c.handlerMutex.Lock()
	if c.handler != nil {
		c.handler(ctx, []byte(data), nil)
	}
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
			status, err := c.client.GetStatus(ctx, nil)
			if err != nil {
				c.reply(ctx, err.Error())
				continue
			}

			if status.Data != "" {
				c.reply(ctx, status.Data)
			}

			time.Sleep(delayBetweenStatusSend)
		}
	}
}
