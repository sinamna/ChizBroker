package broker

import "errors"

var (
	// Use this error for the calls that are coming after that
	// server is shutting down
	ErrUnavailable = errors.New("service is unavailable")
	// Use this error when the message with provided id is not available
	ErrInvalidID = errors.New("message with id provided is not valid or never published")
	// Use this error when message had been published, but it is not
	// available anymore because the expiration time has reached.
	ErrExpiredID = errors.New("message with id provided is expired")
)
