package rpc

import (
	"context"
	"errors"
	"io"
	"net"
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
}

// NewServer creates a new GRPC server.
func NewServer(ctx context.Context, transporter transport.Transporter, netListener net.Listener) (*Server, error) {
	output := Server{
		rpc:         grpc.NewServer(),
		transporter: transporter,
		stop:        make(chan struct{}, 1),
	}

	rpcmodel.RegisterCommandSenderServer(output.rpc, &output)

	reflection.Register(output.rpc)

	transporter.SetResponseHandler(func(ctx context.Context, data []byte, err error) {
		output.bufferMutex.Lock()

		defer output.bufferMutex.Unlock()

		switch {
		case errors.Is(err, io.EOF):
			// Do nothing
		case err != nil:
			output.bufferline = append(output.bufferline, []byte(err.Error()+"\n")...)
		default:
			output.bufferline = append(output.bufferline, data...)
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

	return &rpcmodel.Status{
		Data: string(data),
	}, nil
}
