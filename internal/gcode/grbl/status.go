package grbl

import (
	"os"
	"strconv"
	"strings"

	"github.com/landru29/cnc-serial/internal/model"
)

const (
	positionStart = "MPos:"
)

// CoordinateFromStatus implements the gcode.Processor interface.
func (g Gerbil) CoordinateFromStatus(statusString string) (*model.Status, error) {
	splitter := strings.Split(statusString, positionStart)
	if len(splitter) < 2 {
		return nil, os.ErrNotExist
	}

	match := g.coordinateRegexp.FindAllStringSubmatch(splitter[1], -1)
	if len(match) == 0 {
		return nil, os.ErrNotExist
	}

	if len(match[0]) != 4 {
		return nil, os.ErrNotExist
	}

	xCoordinate, err := strconv.ParseFloat(match[0][1], 64)
	if err != nil {
		return nil, err
	}

	yCoordinate, err := strconv.ParseFloat(match[0][2], 64)
	if err != nil {
		return nil, err
	}

	zCoordinate, err := strconv.ParseFloat(match[0][3], 64)
	if err != nil {
		return nil, err
	}

	return &model.Status{
		XCoordinate: xCoordinate,
		YCoordinate: yCoordinate,
		ZCoordinate: zCoordinate,
	}, nil
}

// BuildStatusRequest implaments the gcode.Processor interface.
func (g Gerbil) BuildStatusRequest() string {
	return "?"
}
