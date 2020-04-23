package client

import (
	"encoding/gob"
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

type Job struct {
	fn   interface{}
	args interface{}
}

type Registry struct {
	store map[string]*Job
}

func (r *Registry) Register(fn interface{}, args interface{}) {
	if args != nil {
		gob.Register(args)
	}
	r.store[runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()] = &Job{
		fn:   fn,
		args: reflect.TypeOf(args),
	}
}

func (r *Registry) Perform(fn string, args interface{}, deadline int) error {
	timeout := time.NewTimer(time.Duration(deadline) * time.Second)
	select {
	case <-timeout.C:
		return timeoutError
	default:
		if args != nil {
			if reflect.TypeOf(args) != reflect.TypeOf(r.store[fn].args) {
				return argsMismatchTypeErr
			}
		}
		registeredFn := reflect.ValueOf(r.store[fn].fn)
		result := registeredFn.Call([]reflect.Value{reflect.ValueOf(r.store[fn].args)})
		resultErr := result[0].Type()

		if resultErr.Implements(errorInterface) {
			resultErrInterface := result[0].Interface()
			return errors.Wrap(resultErrInterface.(error), resultErrInterface.(error).Error())
		}
		return nil
	}
}
