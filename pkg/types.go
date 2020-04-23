package client

import "github.com/pkg/errors"

type Job struct {
	fn   interface{}
	args interface{}
}

type JobErr struct {
	msg       string
	typ       string
	backtrace []string
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
