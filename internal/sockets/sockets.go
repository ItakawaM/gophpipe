package sockets

import (
	"fmt"

	"go.nanomsg.org/mangos/v3/protocol"
	"google.golang.org/protobuf/proto"
)

type SocketFunc func() (protocol.Socket, error)

type Option struct {
	Key   string
	Value any
}

func New(newSock SocketFunc, opts ...Option) (protocol.Socket, error) {
	sock, err := newSock()
	if err != nil {
		return nil, fmt.Errorf("sockets: %w: %w", ErrInvalidSocket, err)
	}

	for _, opt := range opts {
		if err := sock.SetOption(opt.Key, opt.Value); err != nil {
			_ = sock.Close()
			return nil, fmt.Errorf("sockets: %w %q = %v (%T): %w",
				ErrInvalidOption, opt.Key, opt.Value, opt.Value, err)
		}
	}

	return sock, nil
}

func Close(sock protocol.Socket) {
	_ = sock.Close()
}

func SendMessage(sock protocol.Socket, message proto.Message) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("sockets: %w: %w", ErrInvalidMessage, err)
	}

	if err := sock.Send(data); err != nil {
		// we don't have to retry the send now
		// possible errors:
		// ErrSendTimeout, ErrNoPeers, ErrClosed, ErrProtoOp
		return err
	}

	return nil
}

func ReceiveMessage(sock protocol.Socket, dst proto.Message) error {
	data, err := sock.Recv()
	if err != nil {
		// possible errors:
		// ErrRecvTimeout, ErrClosed, ErrProtoOp
		return err
	}

	if err := proto.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("sockets: %w: %w", ErrInvalidMessage, err)
	}

	return nil
}
