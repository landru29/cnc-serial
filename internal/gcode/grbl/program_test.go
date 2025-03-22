package grbl_test

import (
	"os"
	"testing"

	"github.com/landru29/cnc-serial/internal/gcode/grbl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProgram(t *testing.T) {
	file, err := os.Open("testdata/prog01.gcode")
	require.NoError(t, err)

	defer func() {
		_ = file.Close()
	}()

	prog, err := grbl.NewProgram(file)
	require.NoError(t, err)

	for _, expected := range []struct {
		cmd  string
		line int
	}{
		{line: 5, cmd: "G17 G90"},
		{line: 6, cmd: "G21"},
		{line: 9, cmd: "G54"},
		{line: 16, cmd: "M3 S0.0"},
		{line: 21, cmd: "G0 Z5.000"},
		{line: 22, cmd: "G0 X0.000 Y0.000"},
		{line: 23, cmd: "G0 X22.820 Y22.820"},
	} {
		cmd, err := prog.ReadNextInstruction()
		require.NoError(t, err)

		assert.Equal(t, expected.line, prog.CurrentLine())
		require.Equal(t, expected.cmd, cmd)
	}
}
