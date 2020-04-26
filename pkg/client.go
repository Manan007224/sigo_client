package client

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"reflect"
	"sync"

	pb "github.com/Manan007224/sigo/pkg/proto"
	"github.com/Manan007224/sigo_client/pkg/pool"
	"github.com/pkg/errors"
)

var (
	grpcConnUnavailable = errors.New("grpc client conn not avaliable")
)

type Client struct {
	ctx         context.Context
	cancel      context.CancelFunc
	workerWg    *sync.WaitGroup
	Concurrency int
	executor    *Executor
	pool        pool.ConnPool
}

type Worker struct {
	ctx    context.Context
	client *Client
}

func NewWorker(ctx context.Context, client *Client) *Worker {
	return &Worker{ctx: ctx, client: client}
}

func (w *Worker) work() {
	w.client.workerWg.Add(1)
	defer w.client.workerWg.Done()
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
		}

		// fetch job and process it
		job, err := w.Fetch("high")
		if err != nil {
			log.Printf("error in fetching job: %v", err)
			continue
		}

		argType, err := w.client.executor.GetArgType(job.Name)
		if err != nil {
			log.Printf("func not registered %v", err)
			continue
		}
		args := reflect.New(argType).Interface()
		if err = gob.NewDecoder(bytes.NewBuffer(job.Args)).Decode(args); err != nil {
			log.Printf("error in decoding args %v", err)
			continue
		}

		err, jobErr := w.client.executor.Perform(job.Name, args, job.ReserveFor, 10)
		var failPayload *pb.FailPayload
		if err != nil {
			failPayload = &pb.FailPayload{
				Id:           job.Jid,
				ErrorMessage: err.Error(),
				ErrorType:    "ExecutionError",
			}
		} else if jobErr != nil {
			failPayload = &pb.FailPayload{
				Id:           job.Jid,
				ErrorMessage: jobErr.Msg,
				ErrorType:    jobErr.Typ,
				Backtrace:    jobErr.Backtrace,
			}
		}

		if failPayload != nil {
			if err = w.Fail(failPayload); err != nil {
				log.Printf("error in failing job: %v", err)
				continue
			}
		}

		// send ACK message
		if err = w.Acknowledge(job); err != nil {
			log.Printf("error in acking job: %v", err)
		}
	}
}

func (w *Worker) Fetch(queue string) (*pb.JobPayload, error) {
	conn, err := w.client.pool.Get()
	if err != nil {
		return nil, grpcConnUnavailable
	}
	defer w.client.pool.Return(conn)
	sc := pb.NewSchedulerClient(conn)
	return sc.Fetch(context.Background(), &pb.Queue{Name: queue})
}

func (w *Worker) Fail(payload *pb.FailPayload) error {
	conn, err := w.client.pool.Get()
	if err != nil {
		return grpcConnUnavailable
	}
	defer w.client.pool.Return(conn)
	sc := pb.NewSchedulerClient(conn)
	_, err = sc.Fail(context.Background(), payload)
	return err
}

func (w *Worker) Acknowledge(payload *pb.JobPayload) error {
	conn, err := w.client.pool.Get()
	if err != nil {
		return grpcConnUnavailable
	}
	defer w.client.pool.Return(conn)
	sc := pb.NewSchedulerClient(conn)
	_, err = sc.Acknowledge(context.Background(), payload)
	return err
}
