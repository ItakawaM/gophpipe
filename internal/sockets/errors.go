package sockets

import "errors"

var (
	// ErrInvalidSocket is returned when socket creation fails.
	ErrInvalidSocket = errors.New("invalid socket")

	// ErrInvalidOption is returned on an invalid socket option.
	ErrInvalidOption = errors.New("invalid socket option")

	// ErrInvalidMessage is returned on a marshal/unmarshal failure.
	ErrInvalidMessage = errors.New("invalid message")
)
