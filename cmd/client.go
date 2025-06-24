package main

import (
	"context"

	"github.com/landru29/cnc-serial/internal/errors"
	"github.com/landru29/cnc-serial/internal/transport/nop"
	"github.com/landru29/cnc-serial/internal/transport/rpc"
	"github.com/landru29/cnc-serial/internal/transport/serial"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func clientSerialCommand(opts *options) *cobra.Command {
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

			if opts.Serial.PortName != "" && opts.Serial.BitRate > 0 {
				serialClient, err := serial.New(ctx, opts.Serial.PortName, opts.Serial.BitRate)
				if err != nil {
					return err
				}

				app.SetTransport(serialClient)

				return app.Start()
			}

			return errors.ErrNoCommunicationWithMachine
		},
	}

	output.Flags().IntVarP(&opts.Serial.BitRate, "bit-rate", "b", opts.Serial.BitRate, "Bit rate")
	output.Flags().VarP(&opts.Serial.PortName, "port", "p", "Port name")
	output.Flags().Float64VarP(&opts.NavigationInc, "nav-inc", "", opts.NavigationInc, "Navigation increment in millimeters")

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

	output.Flags().Float64VarP(&opts.NavigationInc, "nav-inc", "", opts.NavigationInc, "Navigation increment in millimeters")

	return output
}

func clientRPCCommand(opts *options) *cobra.Command {
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

			if opts.RPC.Clientddr != "" {
				grpcTransport, err := rpc.New(ctx, opts.RPC.Clientddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					return err
				}

				app.SetTransport(grpcTransport)

				return app.Start()
			}

			return errors.ErrNoCommunicationWithMachine
		},
	}

	output.Flags().StringVarP(&opts.RPC.Clientddr, "address", "a", opts.RPC.Clientddr, "RPC server address")
	output.Flags().Float64VarP(&opts.NavigationInc, "nav-inc", "", opts.NavigationInc, "Navigation increment in millimeters")

	return output
}
