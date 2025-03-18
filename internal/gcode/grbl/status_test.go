package grbl_test

import (
	"testing"

	"github.com/landru29/cnc-serial/internal/gcode/grbl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoordinateFromStatus(t *testing.T) {
	ctrl, err := grbl.New()
	require.NoError(t, err)

	status, err := ctrl.CoordinateFromStatus("<Idle,MPos:-4.925,-5.002,-5.000,WPos:-4.925,-5.002,-5.000,Buf:0,RX:1,Ln:0,F:0.>")
	require.NoError(t, err)

	assert.InDelta(t, -4.925, status.XCoordinate, 1e-3)
	assert.InDelta(t, -5.002, status.YCoordinate, 1e-3)
	assert.InDelta(t, -5.000, status.ZCoordinate, 1e-3)
}
