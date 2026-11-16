package sockets

import (
	"testing"
	"testing/synctest"
	"time"

	pb "github.com/ItakawaM/gophpipe/gen/go/testutil/v1"
	"github.com/ItakawaM/gophpipe/internal/testutil"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol"
	"go.nanomsg.org/mangos/v3/protocol/pull"
	"go.nanomsg.org/mangos/v3/protocol/push"
	"google.golang.org/protobuf/proto"
)

// TestNew tests socket creation and options validation.
func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		newSock SocketFunc
		opts    []Option
		err     error
	}{
		{
			name:    "valid socket no options",
			newSock: push.NewSocket,
			opts:    nil,
		},
		{
			name:    "valid socket with options",
			newSock: push.NewSocket,
			opts: []Option{
				{
					Key:   mangos.OptionSendDeadline,
					Value: 1 * time.Second,
				},
			},
		},
		{
			name: "invalid socket",
			newSock: SocketFunc(func() (protocol.Socket, error) {
				return nil, ErrInvalidSocket
			}),
			opts: nil,
			err:  ErrInvalidSocket,
		},
		{
			name:    "invalid options",
			newSock: push.NewSocket,
			opts: []Option{
				{
					Key:   "invalid",
					Value: "invalid",
				},
			},
			err: ErrInvalidOption,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sock, err := New(tt.newSock, tt.opts...)
			if err == nil {
				Close(sock)
			}

			testutil.AssertError(t, err, tt.err)
		})
	}
}

// TestClose tests that the socket is closed correctly.
func TestClose(t *testing.T) {
	sock, err := New(push.NewSocket)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	Close(sock)

	// socket is already closed
	err = sock.Close()
	testutil.AssertError(t, err, mangos.ErrClosed)
}

// TestSendMessage tests message sending and marshaling.
func TestSendMessage(t *testing.T) {
	tests := []struct {
		name        string
		newSock     SocketFunc
		opts        []Option
		message     proto.Message
		wantElapsed time.Duration
		err         error
	}{
		{
			name:    "message queued",
			newSock: push.NewSocket,
			message: &pb.Dummy{
				Data: "123",
			},
		},
		{
			name:    "invalid message",
			newSock: push.NewSocket,
			message: &pb.Dummy{
				Data: string([]byte{0xF0, 0x90, 0x8D}),
			},
			err: ErrInvalidMessage,
		},
		{
			name:    "invalid socket",
			newSock: pull.NewSocket,
			message: &pb.Dummy{
				Data: "123",
			},
			err: mangos.ErrProtoOp,
		},
		{
			name:    "timeout",
			newSock: push.NewSocket,
			opts: []Option{
				{
					Key: mangos.OptionSendDeadline,
					// fake time is used
					Value: 5 * time.Second,
				},
				{
					Key:   mangos.OptionWriteQLen,
					Value: 0,
				},
			},
			message: &pb.Dummy{
				Data: "123",
			},
			wantElapsed: 5 * time.Second,
			err:         mangos.ErrSendTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				sock, err := New(tt.newSock, tt.opts...)
				if err != nil {
					t.Fatalf("New() = %v", err)
				}
				defer Close(sock)

				start := time.Now()

				err = SendMessage(sock, tt.message)
				testutil.AssertError(t, err, tt.err)
				testutil.AssertEqual(t, time.Since(start), tt.wantElapsed)

			})
		})
	}

	// separate case for sending a message on a closed socket
	t.Run("closed socket", func(t *testing.T) {
		sock, err := New(push.NewSocket)
		if err != nil {
			t.Fatalf("New() = %v", err)
		}
		Close(sock)

		err = SendMessage(sock, &pb.Dummy{})
		testutil.AssertError(t, err, mangos.ErrClosed)
	})
}

// TestReceiveMessage tests ReceiveMessage errors.
func TestReceiveMessage(t *testing.T) {
	tests := []struct {
		name        string
		newSock     SocketFunc
		opts        []Option
		wantElapsed time.Duration
		err         error
	}{
		{
			name:    "invalid socket",
			newSock: push.NewSocket,
			opts:    nil,
			err:     mangos.ErrProtoOp,
		},
		{
			name:    "timeout",
			newSock: pull.NewSocket,
			opts: []Option{
				{
					Key:   mangos.OptionRecvDeadline,
					Value: 5 * time.Second,
				},
			},
			wantElapsed: 5 * time.Second,
			err:         mangos.ErrRecvTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				sock, err := New(tt.newSock, tt.opts...)
				if err != nil {
					t.Fatalf("New() = %v", err)
				}
				defer Close(sock)

				start := time.Now()

				var dummy pb.Dummy
				err = ReceiveMessage(sock, &dummy)
				testutil.AssertError(t, err, tt.err)
				testutil.AssertEqual(t, time.Since(start), tt.wantElapsed)
			})
		})
	}

	// separate case for receiving a message on a closed socket
	t.Run("closed socket", func(t *testing.T) {
		sock, err := New(pull.NewSocket)
		if err != nil {
			t.Fatalf("New() = %v", err)
		}
		Close(sock)

		err = ReceiveMessage(sock, &pb.Dummy{})
		testutil.AssertError(t, err, mangos.ErrClosed)
	})
}

// TestSockets tests that ReceiveMessage correctly receives
// the message SendMessage sends.
func TestSockets(t *testing.T) {
	tests := []struct {
		name    string
		message *pb.Dummy
	}{
		{
			name: "full flow",
			message: &pb.Dummy{
				Data: "hello world",
			},
		},
		{
			name:    "nil message",
			message: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender, receiver := newPair(t, push.NewSocket, pull.NewSocket)

			if err := SendMessage(sender, tt.message); err != nil {
				t.Fatalf("SendMessage() = %v", err)
			}

			var out pb.Dummy
			if err := ReceiveMessage(receiver, &out); err != nil {
				t.Fatalf("ReceiveMessage() = %v", err)
			}

			testutil.AssertEqual(t, out.GetData(), tt.message.GetData())
		})
	}

	// separate case for receiving an invalid proto message
	t.Run("invalid message", func(t *testing.T) {
		sender, receiver := newPair(t, push.NewSocket, pull.NewSocket)

		if err := sender.Send([]byte{0xF0, 0x90, 0x8D}); err != nil {
			t.Fatalf("Send() = %v", err)
		}

		var out pb.Dummy
		err := ReceiveMessage(receiver, &out)
		testutil.AssertError(t, err, ErrInvalidMessage)
	})
}
