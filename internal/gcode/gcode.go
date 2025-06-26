// Package gcode manages the g-code.
package gcode

import (
	"github.com/landru29/cnc-serial/internal/lang"
	"github.com/landru29/cnc-serial/internal/model"
)

const (
	DefaultHelperCode = "default"
)

// Processor provides metheds to help the user on G-Codes.
type Processor interface {
	CodeDescription(lang lang.Language, code string) string
	AvailableLanguages() []lang.Language
	Colorize(text string) string
	UnmarshalStatus(statusString string) (*model.Status, error)
	MoveAxis(offset float64, axisName string) string
	CommandStatus() string
	CommandAbsoluteCoordinate() string
	CommandRelativeCoordinate() string
}

// Programmer is the program interface.
type Programmer interface {
	// CurrentLine is the line of the current instruction.
	CurrentLine() int64

	// CurrentCommand is the current selected command.
	CurrentCommand() string

	// Content is the line of the current instruction.
	Content() string

	// Reset reset the program to the first command.
	Reset() error

	// SetLinesToExecute defines the sequence to execute.
	SetLinesToExecute(count int64)

	// ToModel is the data converter.
	ToModel() *model.Program

	// ReadNextInstruction move the pointer to the next instruction and read it.
	ReadNextInstruction() error

	// NextCommandToExecute delevers the next instruction to execute depending on what the user requested.
	NextCommandToExecute() (string, error)
}
