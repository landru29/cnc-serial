// Package main is the main application entrypoint.
package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/landru29/cnc-serial/internal/application"
	"github.com/landru29/cnc-serial/internal/gcode/grbl"
	"github.com/landru29/cnc-serial/internal/lang"
	"github.com/landru29/cnc-serial/internal/stack/memory"
	"github.com/landru29/cnc-serial/internal/transport/nop"
	"github.com/landru29/cnc-serial/internal/transport/serial"
	"github.com/spf13/cobra"
)

const defaultBitRate = 115200

func mainCommand() (*cobra.Command, error) {
	var (
		portName string
		bitRate  int
		dryRun   bool
		language = lang.DefaultLanguage
	)

	stacker := memory.New()

	gerbil, err := grbl.New()
	if err != nil {
		return nil, err
	}

	output := &cobra.Command{
		Use:   "serial",
		Short: "Serial monitor",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer func() {
				cancel()
			}()

			app, err := application.NewClient(ctx, stacker, gerbil)
			if err != nil {
				return err
			}

			app.SetLanguage(language)

			nopTransport := nop.New()

			app.SetTransport(nopTransport)

			if !dryRun {
				serialClient, err := serial.New(portName, bitRate)
				if err != nil {
					return err
				}

				app.SetTransport(serialClient)

				defer func() {
					_ = app.Close()
				}()
			}

			return app.Start()
		},
	}

	output.Flags().IntVarP(&bitRate, "bit-rate", "b", defaultBitRate, "Bit rate")
	output.Flags().StringVarP(&portName, "port", "p", defaultPort(), "Port name")
	output.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run (do not open serial port)")
	output.Flags().VarP(
		&language,
		"lang",
		"l",
		fmt.Sprintf("language (available: %s)", joinLang(gerbil.AvailableLanguages())),
	)

	return output, nil
}

func joinLang(languages []lang.Language) string {
	output := make([]string, len(languages))
	for key, language := range languages {
		output[key] = string(language)
	}

	return strings.Join(output, ", ")
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
