package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/landru29/cnc-serial/internal/errors"
	"github.com/landru29/cnc-serial/internal/transport"
	"github.com/landru29/cnc-serial/internal/transport/nop"
	"github.com/landru29/cnc-serial/internal/transport/rpc"
	"github.com/landru29/cnc-serial/internal/transport/serial"
	"github.com/spf13/cobra"
)

func agentCommand(opts *options) *cobra.Command {
	output := &cobra.Command{
		Use:   "agent",
		Short: "manage the local agent",
	}

	output.AddCommand(rpcAgentCommand(opts))

	return output
}

func rpcAgentCommand(opts *options) *cobra.Command {
	var addr string

	output := &cobra.Command{
		Use:   "rpc",
		Short: "start the rpc agent",
	}

	output.PersistentFlags().StringVarP(&addr, "address", "a", ":1324", "RPC server address")

	output.AddCommand(
		rpcSerialCommand(opts, &addr),
		rpcMockCommand(opts, &addr),
	)

	return output
}

func rpcSerialCommand(opts *options, addr *string) *cobra.Command {
	var (
		bitRate  int
		portName string
	)

	output := &cobra.Command{
		Use:   "serial [filename]",
		Short: "CNC Serial monitor",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var transporter transport.Transporter

			ctx, cancel := context.WithCancel(cmd.Context())
			defer func() {
				cancel()
			}()

			lis, err := net.Listen("tcp", *addr)
			if err != nil {
				return err
			}

			opts.logger.Info("listening gRPC", "addr", *addr)

			if portName != "" && bitRate > 0 {
				serialClient, err := serial.New(ctx, portName, bitRate)
				if err != nil {
					return err
				}

				transporter = serialClient

				opts.logger.Info(fmt.Sprintf("Connected to %s with bitrate %d", portName, bitRate))

				defer func() {
					_ = serialClient.Close()
				}()

				servers, err := rpc.NewServer(ctx, opts.logger, transporter, lis)
				if err != nil {
					return err
				}

				defer func() {
					_ = servers.Close()
				}()

				signalCh := make(chan os.Signal, 1)
				signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

				<-signalCh

				return nil
			}

			return errors.ErrNoCommunicationWithMachine
		},
	}

	output.Flags().IntVarP(&bitRate, "bit-rate", "b", defaultBitRate, "Bit rate")
	output.Flags().StringVarP(&portName, "port", "p", defaultPort(), "Port name")

	return output
}

func rpcMockCommand(opts *options, addr *string) *cobra.Command {
	output := &cobra.Command{
		Use:   "mock [filename]",
		Short: "CNC mock monitor",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var transporter transport.Transporter

			ctx, cancel := context.WithCancel(cmd.Context())
			defer func() {
				cancel()
			}()

			lis, err := net.Listen("tcp", *addr)
			if err != nil {
				return err
			}

			opts.logger.Info("listening gRPC", "addr", *addr)

			nopTransport := nop.New(ctx)

			transporter = nopTransport

			opts.logger.Info("Connected to mock")

			defer func() {
				_ = nopTransport.Close()
			}()

			servers, err := rpc.NewServer(ctx, opts.logger, transporter, lis)
			if err != nil {
				return err
			}

			defer func() {
				_ = servers.Close()
			}()

			signalCh := make(chan os.Signal, 1)
			signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

			<-signalCh

			return nil
		},
	}

	return output
}
