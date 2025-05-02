package main

import (
	"errors"
	"io"
	"os"

	apperrors "github.com/landru29/cnc-serial/internal/errors"
	"github.com/landru29/cnc-serial/internal/gcode/grbl"
	"github.com/spf13/cobra"
)

func cleanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clean <file>",
		Short: "clean gcode",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("missing file")
			}

			file, err := os.Open(args[0])
			if err != nil {
				return err
			}

			defer func(closer io.Closer) {
				_ = closer.Close()
			}(file)

			program, err := grbl.NewProgram(file)
			if err != nil {
				return err
			}

			program.SetLinesToExecute(-1)

			if _, err := program.NextCommandToExecute(); err != nil {
				return err
			}

			var (
				errProg error
				line    = "_"
			)

			for !errors.Is(errProg, apperrors.ErrProgramIdle) && line != "" {
				line, errProg = program.NextCommandToExecute()

				if _, err := cmd.OutOrStdout().Write([]byte(line + "\n")); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
