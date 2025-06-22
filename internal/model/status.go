package model

import (
	"encoding/json"
	"io"
	"strings"
)

const (
	statusMagic            = "status"
	formatedConnectionSize = 20
)

// Status is the CNC status.
type Status struct {
	Machine             *Coordinates     `json:"mpos"`
	Tool                *Coordinates     `json:"wpos"`
	ToolOffset          *Coordinates     `json:"wco"`
	Speed               *Speed           `json:"speed"`
	Alarm               *Alarm           `json:"alarm"`
	State               State            `json:"state"`
	RelativeCoordinates RelativeAbsolute `json:"relativeCoordinates"`
	Type                ObjectName       `json:"type"`
	RemainingProgram    int64            `json:"remainingProgram"`
	CanRun              bool             `json:"canRun"`
	Connection          string           `json:"connection"`
	Buffer              BufferState      `json:"buffer"`
	Overrides           Overrides        `json:"feedOverride"`
}

// FormatedConnection is the connection status.
func (s Status) FormatedConnection() string {
	return (s.Connection + strings.Repeat(" ", formatedConnectionSize))[:formatedConnectionSize]
}

// Encode is the status Encoder.
func (s Status) Encode(writer io.Writer) error {
	s.Type = statusMagic

	return encode(writer, s)
}

// DecodeStatus is the status decoder.
func DecodeStatus(data string) *Status {
	var status Status
	if err := json.Unmarshal([]byte(data), &status); err == nil {
		if status.Type != statusMagic {
			return nil
		}

		return &status
	}

	return nil
}

// ToolCoordinates are the tool coordinate.
func (s Status) ToolCoordinates() Coordinates {
	output := Coordinates{}
	if s.Machine != nil {
		output.XCoordinate += s.Machine.XCoordinate
		output.YCoordinate += s.Machine.YCoordinate
		output.ZCoordinate += s.Machine.ZCoordinate
	}

	if s.ToolOffset != nil {
		output.XCoordinate -= s.ToolOffset.XCoordinate
		output.YCoordinate -= s.ToolOffset.YCoordinate
		output.ZCoordinate -= s.ToolOffset.ZCoordinate
	}

	return output
}

// CurrentState is the formated state.
func (s Status) CurrentState() string {
	output := string(s.State) + "               "

	return strings.ToUpper(output[:8])
}

// State is the CNC state.
type State string

// Alarm is an alarm code.
type Alarm uint32

// Coordinates is 3D coordinates in millimeters.
type Coordinates struct {
	XCoordinate float64 `json:"x"`
	YCoordinate float64 `json:"y"`
	ZCoordinate float64 `json:"z"`
}

func (c *Coordinates) merge(other *Coordinates) *Coordinates {
	switch {
	case c == nil && other == nil:
		return nil
	case c == nil && other != nil:
		return other
	case c != nil && other == nil:
		return c
	default:
		c.XCoordinate = other.XCoordinate
		c.YCoordinate = other.YCoordinate
		c.ZCoordinate = other.ZCoordinate

		return c
	}
}

func (s *Speed) merge(other *Speed) *Speed {
	switch {
	case s == nil && other == nil:
		return nil
	case s == nil && other != nil:
		return other
	case s != nil && other == nil:
		return s
	default:
		s.Spindle = other.Spindle
		s.FeedRate = other.FeedRate

		return s
	}
}

// Speed is feed rate and spindle rotation.
type Speed struct {
	FeedRate float64 `json:"feedRate"` // in millimeters per minute.
	Spindle  float64 `json:"spindle"`  // in rotations per minute.
}

// Merge combines statuses.
func (s *Status) Merge(status Status) {
	s.Machine = s.Machine.merge(status.Machine)
	s.Tool = s.Tool.merge(status.Tool)
	s.ToolOffset = s.ToolOffset.merge(status.ToolOffset)
	s.State = status.State
	s.Speed = s.Speed.merge(status.Speed)

	if s.Alarm == nil || status.Alarm != nil {
		s.Alarm = status.Alarm
	}
}

// RelativeAbsolute is whether coordinates are relatives absolutes.
type RelativeAbsolute bool

// String implements the Stringer interface.
func (r RelativeAbsolute) String() string {
	if r {
		return "REL"
	}

	return "ABS"
}

// BufferState is the status of the buffer.
type BufferState struct {
	AvailableBlocks int64 `json:"availableBloks"`
	AvailableBuffer int64 `json:"availableBuffer"`
}

// Overrides is the defined feeds.
type Overrides struct {
	Feed    int64 `json:"feedX"`
	Rapid   int64 `json:"feedY"`
	Spindle int64 `json:"feedZ"`
}
