// Package model gathers all data models.
package model

// Status is the CNC status.
type Status struct {
	XCoordinate float64 `json:"x"`
	YCoordinate float64 `json:"y"`
	ZCoordinate float64 `json:"z"`
}
