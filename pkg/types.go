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

type stackTracer interface {
	StackTrace() errors.StackTrace
}
