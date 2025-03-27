package rpc

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"strings"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/landru29/cnc-serial/internal/transport"
	rpcmodel "github.com/landru29/cnc-serial/internal/transport/rpc/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server is the RPC server.
type Server struct {
	rpcmodel.UnimplementedCommandSenderServer
	rpc         *grpc.Server
	transporter transport.Transporter
	bufferline  []byte
	bufferMutex sync.Mutex
	stop        chan struct{}
	logger      *slog.Logger
}

// NewServer creates a new GRPC server.
func NewServer(
	ctx context.Context,
	logger *slog.Logger,
	transporter transport.Transporter,
	netListener net.Listener,
) (*Server, error) {
	output := Server{
		rpc:         grpc.NewServer(),
		transporter: transporter,
		stop:        make(chan struct{}, 1),
		logger:      logger,
	}

	rpcmodel.RegisterCommandSenderServer(output.rpc, &output)

	reflection.Register(output.rpc)

	transporter.SetResponseHandler(func(_ context.Context, data []byte, err error) {
		output.bufferMutex.Lock()

		defer output.bufferMutex.Unlock()

		switch {
		case errors.Is(err, io.EOF):
			// Do nothing
		case err != nil:
			output.bufferline = append(output.bufferline, []byte(err.Error()+"\n")...)
			output.logger.Error("handler error", "message", err.Error())

		default:
			output.bufferline = append(output.bufferline, data...)
			output.logger.Info("handler response", "response", strings.ReplaceAll(string(data), "\n", "↲"))
		}
	})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-output.stop:
				return
			default:
				if err := output.rpc.Serve(netListener); err != nil {
					return
				}
			}
		}
	}()

	return &output, nil
}

// SendCommand implements the rpcmodel.CommandSenderServer interface.
func (s *Server) SendCommand(ctx context.Context, cmd *rpcmodel.Command) (*empty.Empty, error) {
	s.logger.Info("Send command", "command", strings.ReplaceAll(cmd.GetData(), "\n", "↲"))

	return nil, s.transporter.Send(ctx, cmd.GetData())
}

// Close implements the io.Closer interface.
func (s *Server) Close() error {
	s.stop <- struct{}{}

	return nil
}

// GetStatus implements the rpcmodel.CommandSenderServer interface.
func (s *Server) GetStatus(context.Context, *empty.Empty) (*rpcmodel.Status, error) {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()

	data := s.bufferline
	s.bufferline = []byte{}

	s.logger.Info("Get status", "status", strings.ReplaceAll(string(data), "\n", "↲"))

	return &rpcmodel.Status{
		Data: string(data),
	}, nil
}
