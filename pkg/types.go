package client

import (
	"reflect"

	"github.com/pkg/errors"
)

type Job struct {
	fn   interface{}
	args reflect.Type
}

type JobErr struct {
	Msg       string
	Typ       string
	Backtrace []string
}

type JobParams struct {
	Fn         string
	Args       interface{}
	Retry      int32
	EnqueueIn  int64
	ReserveFor int64
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
