// Package control is the way to send commands.
package control

import (
	"context"

	"github.com/landru29/cnc-serial/internal/gcode"
	"github.com/landru29/cnc-serial/internal/transport"
)

// Commander is the interface for sending commands.
type Commander interface {
	PushCommands(ctx context.Context, commands ...string) error
	PushProgramCommands(commands ...string) error
	MoveRelative(ctx context.Context, offset float64, axisName string) error
	SetTransporter(transporter transport.Transporter)
	SetProgrammer(programmer gcode.Programmer)
	ProcessResponse(ctx context.Context, data []byte, err error)
}
