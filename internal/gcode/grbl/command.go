package grbl

import (
	"fmt"
	"strings"
)

// CommandStatus implements the gcode.Processor interface.
func (g Gerbil) CommandStatus() string {
	return "?"
}

// CommandAbsoluteCoordinate implements the gcode.Processor interface.
func (g Gerbil) CommandAbsoluteCoordinate() string {
	return "G90"
}

// CommandRelativeCoordinate implements the gcode.Processor interface.
func (g Gerbil) CommandRelativeCoordinate() string {
	return "G91"
}

// MoveAxis implements the gcode.Processor interface.
func (g Gerbil) MoveAxis(offset float64, axisName string) string {
	return fmt.Sprintf("G0 %s%.3f", strings.ToUpper(axisName), offset)
}
