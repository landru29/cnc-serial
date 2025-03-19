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
	BuildStatusRequest() string
	UnmarshalStatus(statusString string) (*model.Status, error)
}
