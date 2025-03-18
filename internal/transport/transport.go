// Package transport is the command transportation to the CNC.
package transport

import "io"

// Transporter is the data transport.
type Transporter interface {
	io.Reader
	Send(commands ...string) error
}

// TransportCloser is a Transporter that can be closed.
type TransportCloser interface {
	io.Closer
	Transporter
}
