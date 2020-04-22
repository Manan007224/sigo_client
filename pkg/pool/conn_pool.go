package pool

import (
	"log"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	addr         = "localhost:50051"
	timeOutError = errors.New("timeout waiting to build connection")
	closedPool   = errors.New("the pool is closed")
)

// Most of the ideas and code taken from this blog - http://dustin.sallings.org/2014/04/25/chan-pool.html
type ConnPool struct {
	FreeConns     chan *grpc.ClientConn
	UsedConns     chan struct{}
	CreateHandler func() (*grpc.ClientConn, error)
	CloseHandler  func(conn *grpc.ClientConn)
}

func NewConnPool(poolSize int, poolOverflow int) *ConnPool {
	return &ConnPool{
		FreeConns:     make(chan *grpc.ClientConn, poolSize),
		UsedConns:     make(chan struct{}, poolOverflow),
		CreateHandler: createClientConn,
		CloseHandler:  closeClientConn,
	}
}

func createClientConn() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create grpc client conn")
	}
	return conn, err
}

func closeClientConn(conn *grpc.ClientConn) {
	if err := conn.Close(); err != nil {
		log.Printf("Error in closing client conn: %v", err)
	}
}

func (cp *ConnPool) Get() (*grpc.ClientConn, error) {
	if cp == nil {
		return nil, errors.New("no pool")
	}

	timeout := time.NewTimer(time.Millisecond)
	defer timeout.Stop()

	select {
	case conn, ok := <-cp.FreeConns:
		if !ok {
			return nil, closedPool
		}
		return conn, nil
	case <-timeout.C:
		timeout.Reset(time.Millisecond)
		select {
		case conn, ok := <-cp.FreeConns:
			if !ok {
				return nil, closedPool
			}
			return conn, nil
		case cp.UsedConns <- struct{}{}:
			conn, err := cp.CreateHandler()
			if err != nil {
				<-cp.UsedConns
				return nil, err
			}
			return conn, nil
		case <-timeout.C:
			return nil, timeOutError
		}
	}
}

func (cp *ConnPool) Return(conn *grpc.ClientConn) {
	defer func() {
		if recover() != nil {
			conn.Close()
		}
	}()

	select {
	case cp.FreeConns <- conn:
	default:
		<-cp.UsedConns
		cp.CloseHandler(conn)
	}
}

func (cp *ConnPool) Close() (err error) {
	defer func() {
		err, _ = recover().(error)
	}()
	close(cp.FreeConns)
	for conn := range cp.FreeConns {
		conn.Close()
	}
	return
}
