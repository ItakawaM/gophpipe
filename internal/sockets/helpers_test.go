package sockets

import (
	"testing"
	"time"

	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol"
	_ "go.nanomsg.org/mangos/v3/transport/inproc"
)

// newPair creates a sender and receiver of the provided protocol.
// It uses inproc as transport and can be safely used in parallel tests.
//
// It closes the sockets in a Cleanup.
func newPair(t *testing.T, newSender, newReceiver SocketFunc) (sender, receiver protocol.Socket) {
	t.Helper()

	addr := "inproc://" + t.Name()

	receiver, err := New(newReceiver, Option{
		Key:   mangos.OptionRecvDeadline,
		Value: time.Second,
	})
	if err != nil {
		t.Fatalf("New(receiver) = %v", err)
	}
	t.Cleanup(func() {
		Close(receiver)
	})

	if err := receiver.Listen(addr); err != nil {
		t.Fatalf("Listen(%q) = %v", addr, err)
	}

	sender, err = New(newSender)
	if err != nil {
		t.Fatalf("New(sender) = %v", err)
	}
	t.Cleanup(func() {
		Close(sender)
	})

	if err := sender.Dial(addr); err != nil {
		t.Fatalf("Dial(%q) = %v", addr, err)
	}

	return sender, receiver
}
