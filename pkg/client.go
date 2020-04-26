package client

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/Manan007224/sigo_client/pkg/pool"

	pb "github.com/Manan007224/sigo/pkg/proto"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	grpcConnUnavailable = errors.New("grpc client conn not avaliable")
	queueChoices        = []Choice{
		{Item: "High", Weight: 4},
		{Item: "Medium", Weight: 2},
		{Item: "Low", Weight: 1},
	}
)

type Client struct {
	ctx          context.Context
	cancel       context.CancelFunc
	workerWg     *sync.WaitGroup
	hearbeatWg   *sync.WaitGroup
	disconnected chan struct{}
	jobs         chan *pb.JobPayload
	Concurrency  int
	executor     *Executor
	pool         *pool.ConnPool
	queuePicker  *Chooser
}

func NewClient(concurrency int) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		ctx:         ctx,
		cancel:      cancel,
		workerWg:    &sync.WaitGroup{},
		hearbeatWg:  &sync.WaitGroup{},
		Concurrency: concurrency,
		executor:    NewExecutor(),
		pool:        pool.NewConnPool(concurrency, concurrency),
		queuePicker: NewChooser([]Choice{}),
	}
}

func (c *Client) Run() {
	config := &pb.ClientConfig{
		Id: "123",
		// TODO - use guid package to generate unique id's
		Type: "golang",
	}
	for _, qc := range queueChoices {
		c.queuePicker.AddChoice(qc)
		config.Queues = append(config.Queues, &pb.QueueConfig{
			Name:     qc.Item.(string),
			Priority: int32(qc.Weight),
		})
	}

	if _, err := c.Query("discover", config); err != nil {
		log.Printf("error in discover client %v", err)
		return
	}

	go c.hearbeat()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-c.disconnected:
		c.gracefullShutdown()
		return
	case <-sigs:
		c.gracefullShutdown()
		return
	}
}

func (c *Client) Perform(params *JobParams) error {
	if !c.executor.Validate(params.Fn) {
		return funcNotRegistered
	}

	b := new(bytes.Buffer)
	defer b.Reset()
	if err := gob.NewEncoder(b).Encode(params.Args); err != nil {
		return err
	}

	jobPayload := &pb.JobPayload{
		Name:       params.Fn,
		Args:       b.Bytes(),
		ReserveFor: params.ReserveFor,
		Retry:      params.Retry,
		Jid:        uuid.NewV4().String(),
	}
	if params.EnqueueIn != 0 {
		jobPayload.EnqueueAt = time.Now().Add(time.Duration(params.EnqueueIn)).Unix()
	}
	if _, err := c.Query("broadcast", jobPayload); err != nil {
		log.Printf("error in broadcasting job: %v", err)
		return err
	}
	return nil
}

func (c *Client) gracefullShutdown() {
	c.cancel()
	c.hearbeatWg.Wait()
	c.workerWg.Wait()
	c.pool.Close()
}

func (c *Client) hearbeat() {
	c.hearbeatWg.Add(1)
	ticker := time.NewTicker(10)
	firstRun := false
	for {
		select {
		case <-c.ctx.Done():
			ticker.Stop()
			c.hearbeatWg.Done()
			return
		case <-ticker.C:
			if _, err := c.Query("beat", &emptypb.Empty{}); err != nil {
				log.Printf("error in heartbeat %v", err)
				if firstRun {
					c.disconnected <- struct{}{}
					return
				} else {
					firstRun = false
				}
			}
		}
	}
}

func (c *Client) Query(rpc string, args interface{}) (interface{}, error) {
	conn, err := c.pool.Get()
	if err != nil {
		return nil, grpcConnUnavailable
	}
	defer c.pool.Return(conn)
	schedulerClient := pb.NewSchedulerClient(conn)

	ctx := context.Background()

	switch rpc {
	case "fail":
		return schedulerClient.Fail(ctx, args.(*pb.FailPayload))
	case "acknowledge":
		return schedulerClient.Acknowledge(ctx, args.(*pb.JobPayload))
	case "fetch":
		return schedulerClient.Fetch(ctx, args.(*pb.Queue))
	case "discover":
		return schedulerClient.Discover(ctx, args.(*pb.ClientConfig))
	case "beat":
		return schedulerClient.HeartBeat(ctx, args.(*emptypb.Empty))
	default:
		return schedulerClient.BroadCast(ctx, args.(*pb.JobPayload))
	}
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
		payload, err := w.client.Query("fetch", w.client.queuePicker.Pick().(string))
		job := payload.(*pb.JobPayload)
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
			if _, err = w.client.Query("fail", failPayload); err != nil {
				log.Printf("error in failing job: %v", err)
				continue
			}
		}

		// send ACK message
		if _, err = w.client.Query("acknowledge", job); err != nil {
			log.Printf("error in acking job: %v", err)
		}
	}
}
