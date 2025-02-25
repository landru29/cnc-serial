// Package main is the main application entrypoint.
package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/landru29/serial/internal/application"
	"github.com/landru29/serial/internal/gcode"
	"github.com/spf13/cobra"
)

const defaultBitRate = 115200

func mainCommand() (*cobra.Command, error) {
	var (
		portName string
		bitRate  int
		dryRun   bool
		language = gcode.DefaultLanguage
	)

	app, err := application.NewClient()
	if err != nil {
		return nil, err
	}

	output := &cobra.Command{
		Use:   "serial",
		Short: "Serial monitor",
		RunE: func(_ *cobra.Command, _ []string) error {
			app.SetLanguage(language)

			if !dryRun {
				if err := app.OpenPort(portName, bitRate); err != nil {
					return err
				}

				defer func() {
					_ = app.Close()
				}()

				go app.Bind()
			}

			return app.Start()
		},
	}

	output.Flags().IntVarP(&bitRate, "bit-rate", "b", defaultBitRate, "Bit rate")
	output.Flags().StringVarP(&portName, "port", "p", app.DefaultPort(), "Port name")
	output.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run (do not open serial port)")
	output.Flags().VarP(
		&language,
		"lang",
		"l",
		fmt.Sprintf("language (available: %s)", strings.Join(app.AvailableLanguages(), ", ")),
	)

	return output, nil
}

func main() {
	cmd, err := mainCommand()
	if err != nil {
		panic(err)
	}

	if err := cmd.ExecuteContext(context.Background()); err != nil {
		panic(err)
	}
}
