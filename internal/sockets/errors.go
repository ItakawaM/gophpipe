package sockets

import "errors"

var (
	ErrInvalidSocket = errors.New("invalid socket")

	ErrInvalidOption = errors.New("invalid socket option")

	ErrInvalidMessage = errors.New("invalid message")
)
