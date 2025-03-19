package grbl_test

import (
	"testing"

	"github.com/landru29/cnc-serial/internal/gcode/grbl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoordinateFromStatus(t *testing.T) {
	t.Run("machine position", func(t *testing.T) {
		ctrl, err := grbl.New()
		require.NoError(t, err)

		status, err := ctrl.UnmarshalStatus("<Idle|MPos:-4.925,-5.002,-5.000|WPos:0.000,0.000,0.000|FS:0,0|Pn:P>")
		require.NoError(t, err)

		assert.InDelta(t, -4.925, status.Machine.XCoordinate, 1e-3)
		assert.InDelta(t, -5.002, status.Machine.YCoordinate, 1e-3)
		assert.InDelta(t, -5.000, status.Machine.ZCoordinate, 1e-3)
	})

	t.Run("tool position", func(t *testing.T) {
		ctrl, err := grbl.New()
		require.NoError(t, err)

		status, err := ctrl.UnmarshalStatus("<Idle|MPos:0.000,0.000,0.000|WPos:-4.925,-5.002,-5.000|FS:0,0|Pn:P>")
		require.NoError(t, err)

		assert.InDelta(t, -4.925, status.Tool.XCoordinate, 1e-3)
		assert.InDelta(t, -5.002, status.Tool.YCoordinate, 1e-3)
		assert.InDelta(t, -5.000, status.Tool.ZCoordinate, 1e-3)
	})
	t.Run("speeds", func(t *testing.T) {
		ctrl, err := grbl.New()
		require.NoError(t, err)

		status, err := ctrl.UnmarshalStatus("<Idle|MPos:0.000,0.000,0.000|WPos:-4.925,-5.002,-5.000|FS:15,42|Pn:P>")
		require.NoError(t, err)

		assert.InDelta(t, 15, status.FeedRate, 1e-3)
		assert.InDelta(t, 42, status.SpindleSpeed, 1e-3)
	})

	t.Run("tool offset", func(t *testing.T) {
		ctrl, err := grbl.New()
		require.NoError(t, err)

		status, err := ctrl.UnmarshalStatus("<Idle|MPos:0.000,0.000,0.000|WPos:0.000,0.000,0.000|WCO:-4.925,-5.002,-5.000|FS:0,0|Pn:P>")
		require.NoError(t, err)

		assert.InDelta(t, -4.925, status.ToolOffset.XCoordinate, 1e-3)
		assert.InDelta(t, -5.002, status.ToolOffset.YCoordinate, 1e-3)
		assert.InDelta(t, -5.000, status.ToolOffset.ZCoordinate, 1e-3)
	})
}
