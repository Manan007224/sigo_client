package client

import (
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	addr         = "localhost:50051"
	timeOutError = errors.New("timeout waiting to build connection")
	closedPool   = errors.New("the pool is closed")
)

type ConnPool struct {
	freeConns chan *grpc.ClientConn
	usedConns chan struct{}
}

func NewConnPool(poolSize int) *ConnPool {
	return &ConnPool{
		freeConns: make(chan *grpc.ClientConn, poolSize),
		usedConns: make(chan struct{}, poolSize),
	}
}

func (cp *ConnPool) Get() (*grpc.ClientConn, error) {
	if cp == nil {
		return nil, errors.New("no pool")
	}

	timeout := time.NewTimer(time.Millisecond)
	defer timeout.Stop()

	select {
	case conn, ok := <-cp.freeConns:
		if !ok {
			return nil, closedPool
		}
		return conn, nil
	case <-timeout.C:
		timeout.Reset(time.Millisecond)
		select {
		case conn, ok := <-cp.freeConns:
			if !ok {
				return nil, closedPool
			}
			return conn, nil
		case cp.usedConns <- struct{}{}:
			conn, err := grpc.Dial(addr)
			if err != nil {
				<-cp.usedConns
				return nil, errors.Wrap(err, "failed to create grpc client conn")
			}
			return conn, nil
		}
	case <-timeout.C:
		return nil, timeOutError
	}
}
