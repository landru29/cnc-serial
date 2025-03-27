package grbl

import (
	"os"
	"strconv"
	"strings"

	"github.com/landru29/cnc-serial/internal/model"
)

// UnmarshalStatus implements the gcode.Processor interface.
func (g Gerbil) UnmarshalStatus(statusString string) (*model.Status, error) {
	var output model.Status

	statusString = strings.TrimSpace(statusString)

	if len(statusString) < 5 || statusString[0] != '<' || statusString[len(statusString)-1] != '>' {
		return nil, os.ErrNotExist
	}

	statusString = statusString[1 : len(statusString)-1]

	statusMatch := g.stateRegexp.FindAllStringSubmatch(statusString, -1)
	if len(statusMatch) < 1 || len(statusMatch[0]) < 2 {
		return nil, os.ErrNotExist
	}

	output.State = model.State(statusMatch[0][1])

	statusString = statusString[len(output.State):]

	for statusString != "" {
		statusString = g.decodeNext(statusString, &output)
	}

	return &output, nil
}

func (g Gerbil) decodeNext(statusString string, output *model.Status) string { //nolint: funlen,gocognit,cyclop,gocyclo
	// go to the first word.
	idx := 0
	for idx < len(statusString) &&
		(statusString[idx] < 'A' || statusString[idx] > 'Z') &&
		(statusString[idx] < 'a' || statusString[idx] > 'z') {
		idx++
	}

	statusString = statusString[idx:]

	// Read argument name.
	argMatch := g.argumentRegexp.FindAllStringSubmatch(statusString, -1)
	if len(argMatch) == 0 || len(argMatch[0]) < 2 {
		return ""
	}

	statusString = statusString[len(argMatch[0][1])+1:]

	var err error

	// process arguments.
	switch strings.ToLower(argMatch[0][1]) {
	case "mpos":
		output.Machine = &model.Coordinates{}

		if statusString, output.Machine.XCoordinate, err = g.readNumber(statusString); err != nil {
			return ""
		}

		if statusString, output.Machine.YCoordinate, err = g.readNumber(statusString); err != nil {
			return ""
		}

		if statusString, output.Machine.ZCoordinate, err = g.readNumber(statusString); err != nil {
			return ""
		}

	case "wpos":
		output.Tool = &model.Coordinates{}

		if statusString, output.Tool.XCoordinate, err = g.readNumber(statusString); err != nil {
			return ""
		}

		if statusString, output.Tool.YCoordinate, err = g.readNumber(statusString); err != nil {
			return ""
		}

		if statusString, output.Tool.ZCoordinate, err = g.readNumber(statusString); err != nil {
			return ""
		}

	case "fs":
		output.Speed = &model.Speed{}

		if statusString, output.Speed.FeedRate, err = g.readNumber(statusString); err != nil {
			return ""
		}

		if statusString, output.Speed.Spindle, err = g.readNumber(statusString); err != nil {
			return ""
		}

	case "wco":
		output.ToolOffset = &model.Coordinates{}

		if statusString, output.ToolOffset.XCoordinate, err = g.readNumber(statusString); err != nil {
			return ""
		}

		if statusString, output.ToolOffset.YCoordinate, err = g.readNumber(statusString); err != nil {
			return ""
		}

		if statusString, output.ToolOffset.ZCoordinate, err = g.readNumber(statusString); err != nil {
			return ""
		}
	case "alarm":
		var alarm float64

		if statusString, alarm, err = g.readNumber(statusString); err != nil {
			return ""
		}

		casted := model.Alarm(alarm)

		output.Alarm = &casted

	case "pn":
		for len(statusString) > 0 &&
			((statusString[0] > 'a' && statusString[0] < 'z') ||
				(statusString[0] > 'A' &&
					statusString[0] < 'Z')) {
			statusString = statusString[1:]
		}

	default:
		return ""
	}

	return statusString
}

func (g Gerbil) readNumber(statusString string) (string, float64, error) {
	if len(statusString) == 0 {
		return statusString, 0, os.ErrNotExist
	}

	match := g.numberRegexp.FindAllStringSubmatch(statusString, -1)
	if len(match) == 0 || len(match[0]) < 2 {
		return g.readNumber(statusString[1:])
	}

	out, err := strconv.ParseFloat(match[0][1], 64)

	return statusString[len(match[0][1]):], out, err
}
