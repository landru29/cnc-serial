// Package gcode manages the g-code.
package gcode

import (
	"github.com/landru29/cnc-serial/internal/lang"
	"github.com/landru29/cnc-serial/internal/model"
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
	CurrentLine() int
	CurrentCommand() string
	ReadNextInstruction() (string, error)
	Content() string
	ToModel() *model.Program
}
