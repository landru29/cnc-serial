// Package control is the way to send commands.
package control

// Commander is the interface for sending commands.
type Commander interface {
	PushCommands(commands ...string) error
	IsRelative() bool
}
