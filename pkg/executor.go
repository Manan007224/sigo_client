package client

import (
	"encoding/gob"
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/pkg/errors"
)

var (
	argsMismatchTypeErr = errors.New("arguments mistmatch")
	timeoutError        = errors.New("job execution timeout")
	errorInterface      = reflect.TypeOf((*error)(nil)).Elem()
)

type Executor struct {
	store map[string]*Job
}

func (r *Executor) Register(fn interface{}, args interface{}) {
	if args != nil {
		gob.Register(args)
	}
	r.store[runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()] = &Job{
		fn:   fn,
		args: reflect.TypeOf(args),
	}
}

func (r *Executor) Perform(fn string, args interface{}, deadline int, limit int) (error, *JobErr) {
	timeout := time.NewTimer(time.Duration(deadline) * time.Second)
	select {
	case <-timeout.C:
		return timeoutError, nil
	default:
		if args != nil {
			if reflect.TypeOf(args) != reflect.TypeOf(r.store[fn].args) {
				return argsMismatchTypeErr, nil
			}
		}
		registeredFn := reflect.ValueOf(r.store[fn].fn)
		result := registeredFn.Call([]reflect.Value{reflect.ValueOf(r.store[fn].args)})
		resultErr := result[0].Type()

		if resultErr.Implements(errorInterface) {
			return nil, r.buildJobErr(result[0].Interface().(error), fn, limit)
		}
		return nil, nil
	}
}

func (r *Executor) buildJobErr(err error, fn string, limit int) *JobErr {
	cerr := errors.Cause(err).(stackTracer)
	sterr := cerr.StackTrace()
	return &JobErr{
		msg:       err.Error(),
		typ:       fmt.Sprintf("%sErr", fn),
		backtrace: fmt.Sprintf("%+v", sterr[0:limit]),
	}
}
