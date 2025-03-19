// Package model gathers all data models.
package model

// Status is the CNC status.
type Status struct {
	Machine      Coordinates `json:"mpos"`
	Tool         Coordinates `json:"wpos"`
	ToolOffset   Coordinates `json:"wco"`
	FeedRate     float64     `json:"feedRate"`     // in millimeters per minute.
	SpindleSpeed float64     `json:"spindleSpeed"` // in rotations per minute.
	Alarm        Alarm       `json:"alarm"`
	State        State       `json:"state"`
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
