package client

import (
	"encoding/gob"
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/ztrue/tracerr"
)

var (
	argsMismatchTypeErr = errors.New("arguments mistmatch")
	funcNotRegistered   = errors.New("function not registered")
	timeoutError        = errors.New("job execution timeout")
	errorInterface      = reflect.TypeOf((*error)(nil)).Elem()
)

type Executor struct {
	store map[string]*Job
}

func NewExecutor() *Executor {
	return &Executor{
		store: make(map[string]*Job),
	}
}

func (r *Executor) Register(fn interface{}, fnName string, args interface{}) {
	if args != nil {
		gob.Register(args)
	}
	r.store[fnName] = &Job{
		fn:   fn,
		args: reflect.TypeOf(args),
	}
	fmt.Println(r.store)
}

func (r *Executor) Perform(fn string, args interface{}, deadline int64, limit int) (error, *JobErr) {
	timeout := time.NewTimer(time.Duration(deadline) * time.Second)
	select {
	case <-timeout.C:
		timeout.Stop()
		return timeoutError, nil
	default:
		if _, ok := r.store[fn]; !ok {
			return funcNotRegistered, nil
		}
		if args != nil {
			if reflect.TypeOf(args) != r.store[fn].args {
				return argsMismatchTypeErr, nil
			}
		}
		registeredFn := reflect.ValueOf(r.store[fn].fn)
		result := registeredFn.Call([]reflect.Value{reflect.ValueOf(args)})
		resultErr := result[0].Type()

		if resultErr.Implements(errorInterface) {
			return nil, r.buildJobErr(result[0].Interface().(error), fn, limit)
		}
		return nil, nil
	}
}

func (r *Executor) GetArgType(fn string) (reflect.Type, error) {
	if _, ok := r.store[fn]; ok {
		return nil, funcNotRegistered
	}
	return r.store[fn].args, nil
}

func (r *Executor) buildJobErr(err error, fn string, limit int) *JobErr {
	frames := tracerr.StackTrace(err)
	if len(frames) >= limit {
		frames = frames[0:limit]
	}
	backtrace := []string{}
	for _, frame := range frames {
		backtrace = append(backtrace, frame.String())
	}
	return &JobErr{
		Msg:       err.Error(),
		Typ:       fmt.Sprintf("%sErr", fn),
		Backtrace: backtrace,
	}
}
