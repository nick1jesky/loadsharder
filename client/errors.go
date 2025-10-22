package client

import "errors"

var (
	ErrInvalidParam      = errors.New("invalid param")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrQueueFull         = errors.New("queue is full")
	ErrClientStopped     = errors.New("client is stopped")
)
