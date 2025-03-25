// Package main is the main application entrypoint.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/landru29/cnc-serial/internal/application"
	"github.com/landru29/cnc-serial/internal/gcode/grbl"
	"github.com/landru29/cnc-serial/internal/lang"
	"github.com/landru29/cnc-serial/internal/stack/memory"
	"github.com/spf13/cobra"
)

const defaultBitRate = 115200

type options struct {
	availableLanguages []lang.Language
	language           lang.Language
	gerbil             *grbl.Gerbil
	stacker            *memory.Stack
}

func foo(ctx context.Context, opts *options, args []string) (*application.Client, error) {
	var program *grbl.Program

	if len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			return nil, err
		}

		defer func(closer io.Closer) {
			_ = closer.Close()
		}(file)

		program, err = grbl.NewProgram(file)
		if err != nil {
			return nil, err
		}
	}

	app, err := application.NewClient(ctx, opts.stacker, opts.gerbil, program)
	if err != nil {
		return nil, err
	}

	app.SetLanguage(opts.language)

	return app, nil
}

func mainCommand() (*cobra.Command, error) {
	var (
		opts      options
		forceGRPC bool
	)

	gerbil, err := grbl.New()
	if err != nil {
		return nil, err
	}

	opts.availableLanguages = gerbil.AvailableLanguages()
	opts.gerbil = gerbil
	opts.language = lang.DefaultLanguage
	opts.stacker = memory.New()

	output := &cobra.Command{
		Use:   "cnc-serial [filename]",
		Short: "CNC Serial monitor",
	}

	output.Flags().BoolVarP(&forceGRPC, "grpc", "", false, "RPC connection")
	output.PersistentFlags().VarP(
		&opts.language,
		"lang",
		"l",
		fmt.Sprintf("language (available: %s)", joinLang(opts.availableLanguages)),
	)

	output.AddCommand(
		agentCommand(),
		clientSerialCommand(&opts),
		clientMockCommand(&opts),
		clientRPCCommand(&opts),
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalCh

		cancel()
	}()

	if err := cmd.ExecuteContext(ctx); err != nil {
		panic(err)
	}
}
