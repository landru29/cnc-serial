package main

import (
	"context"
	"errors"

	"github.com/landru29/cnc-serial/internal/transport/nop"
	"github.com/landru29/cnc-serial/internal/transport/rpc"
	"github.com/landru29/cnc-serial/internal/transport/serial"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func clientSerialCommand(opts *options) *cobra.Command {
	var (
		bitRate  int
		portName string
	)

	output := &cobra.Command{
		Use:   "serial [filename]",
		Short: "CNC Serial monitor",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer func() {
				cancel()
			}()

			app, err := initApp(ctx, opts, args)
			if err != nil {
				return err
			}

			defer func() {
				_ = app.Close()
			}()

			if portName != "" && bitRate > 0 {
				serialClient, err := serial.New(ctx, portName, bitRate)
				if err != nil {
					return err
				}

				app.SetTransport(serialClient)

				return app.Start()
			}

			return errors.New("no communication with the machine")
		},
	}

	output.PersistentFlags().IntVarP(&bitRate, "bit-rate", "b", defaultBitRate, "Bit rate")
	output.PersistentFlags().StringVarP(&portName, "port", "p", defaultPort(), "Port name")

	return output
}

func clientMockCommand(opts *options) *cobra.Command {
	output := &cobra.Command{
		Use:   "mock [filename]",
		Short: "CNC mock monitor",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer func() {
				cancel()
			}()

			app, err := initApp(ctx, opts, args)
			if err != nil {
				return err
			}

			defer func() {
				_ = app.Close()
			}()

			nopTransport := nop.New(ctx)

			app.SetTransport(nopTransport)

			return app.Start()
		},
	}

	return output
}

func clientRPCCommand(opts *options) *cobra.Command {
	var addr string

	output := &cobra.Command{
		Use:   "rpc [filename]",
		Short: "CNC RPC monitor",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer func() {
				cancel()
			}()

			app, err := initApp(ctx, opts, args)
			if err != nil {
				return err
			}

			defer func() {
				_ = app.Close()
			}()

			if addr != "" {
				grpcTransport, err := rpc.New(ctx, opts.logger, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					return err
				}

				app.SetTransport(grpcTransport)

				return app.Start()
			}

			return errors.New("no communication with the machine")
		},
	}

	output.Flags().StringVarP(&addr, "address", "a", "0.0.0.0:1324", "RPC server address")

	return output
}
