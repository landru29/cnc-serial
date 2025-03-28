// Package errors manages application errors
package errors

const (
	// ErrNoCommunicationWithMachine is when no communication could be set with the machine.
	ErrNoCommunicationWithMachine Error = "no communication with the machine"

	// ErrMissingTransporter is when no transporter is specified.
	ErrMissingTransporter Error = "missing transporter"

	// ErrUnrecognizedStatus is when a status could not be decoded from a string.
	ErrUnrecognizedStatus Error = "unrecognized status"

	// ErrProgramIdle is when no more instruction can be sent to the machine due to a stop in the program.
	ErrProgramIdle Error = "program is idle"
)

// Error is a string error.
type Error string

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}
