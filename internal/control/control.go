// Package control is the way to send commands.
package control

import (
	"github.com/landru29/cnc-serial/internal/gcode"
	"github.com/landru29/cnc-serial/internal/transport"
)

// Commander is the interface for sending commands.
type Commander interface {
	PushCommands(commands ...string) error
	PushProgramCommands(commands ...string) error
	MoveRelative(offset float64, axisName string) error
	SetTransporter(transporter transport.Transporter)
	SetProgrammer(programmer gcode.Programmer)
}
