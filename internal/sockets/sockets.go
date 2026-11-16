// Package sockets integrates proto into mangos sockets.
package sockets

import (
	"fmt"

	"go.nanomsg.org/mangos/v3/protocol"
	"google.golang.org/protobuf/proto"
)

// SocketFunc is a constructor function for protocol.Socket.
type SocketFunc func() (protocol.Socket, error)

// Option is a single mangos socket option.
type Option struct {
	Key   string
	Value any
}

// New creates a new specified socket and sets its options.
//
// It returns ErrInvalidSocket when socket creation fails (Does not happen as of mangos v3.6.0).
// It returns ErrInvalidOption on an invalid option in opts.
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

// Close closes the socket.
// Exists solely for short defer statements.
//
// It should be called only once.
func Close(sock protocol.Socket) {
	_ = sock.Close()
}

// SendMessage sends/queues a typed message with the socket.
//
// It returns ErrInvalidMessage on a marshal error (example: invalid utf-8 string).
// It returns mangos.ErrProtoOp if the socket cannot Send.
// It returns mangos.ErrSendTimeout on a timeout (example: PUSH full queue and short deadline).
// It returns mangos.ErrClosed if the socket is closed.
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

// ReceiveMessage receives a typed message and unmarshals it into dst.
//
// It returns ErrInvalidMessage on an unmarshal error (example: invalid utf-8 string).
// It returns mangos.ErrProtoOp if the socket cannot Recv.
// It returns mangos.ErrSendTimeout on a timeout.
// It returns mangos.ErrClosed if the socket is closed.
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
