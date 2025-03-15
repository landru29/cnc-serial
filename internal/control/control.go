// Package control is the way to send commands.
package control

import (
	"io"

	"github.com/landru29/serial/internal/stack"
)

// Commander is the interface for sending commands.
type Commander interface {
	io.Closer
	Send(texts ...string) error
	IsRelative() bool
	CommandStack() *stack.Stack
	Bind()
}
