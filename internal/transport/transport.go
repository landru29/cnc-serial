// Package transport is the command transportation to the CNC.
package transport

import (
	"context"
	"io"
)

type ResponseHandler func(context.Context, []byte, error)

// Transporter is the data transport.
type Transporter interface {
	Send(ctx context.Context, commands ...string) error
	SetResponseHandler(ResponseHandler)
}

// TransportCloser is a Transporter that can be closed.
type TransportCloser interface { //nolint: revive
	io.Closer
	Transporter
}
