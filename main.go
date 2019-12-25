package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"time"
	"fmt"
	"encoding/json"
	"bytes"
	"log"
	"github.com/emirpasic/gods/sets/hashset"
	"strconv"
)

var server = "localhost:3000"

type Manager struct {
	pool []*Worker
	workers int
	conn *websocket.Conn
	connectionAlive chan bool
	lastHeartBeat time.Time
	busyWorkers *hashset.Set
	freeWorkers *hashset.Set
}

type Worker struct {
	jobChan chan Job
}


type Job struct {
	Jid string
	Name string
	Args map[string]interface{}
	Queue string
}

func NewManager(concurrency int) (*Manager, error) {
	mgr:= &Manager {
		pool: []*Worker{},
		workers: concurrency,
		connectionAlive: make(chan bool),
		lastHeartBeat: time.Now(),
		freeWorkers: hashset.New(),
		busyWorkers: hashset.New(),
	}

	uri := "ws://localhost:3000/consume"
	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		return nil, fmt.Errorf("error in connectiing to the server: %s", err)
	}

	mgr.conn = conn

	for w := 0; w < mgr.workers; w++ {
		worker := &Worker {
			jobChan: make(chan Job),
		}
		mgr.pool = append(mgr.pool, worker)
		mgr.freeWorkers.Add(worker)
	}
	return mgr, nil
}

func (mgr *Manager) Process() {
	go mgr.heartbeat_server()

	select {
	case <-mgr.connectionAlive:
		fmt.Println("worker is dead")
	}
}

func (mgr *Manager) worker_utilization() string {
	return fmt.Sprintf("utilization-%s", strconv.Itoa(mgr.busyWorkers.Size() / len(mgr.pool)))
}


func (mgr *Manager) heartbeat_server() {
	for {
		latestHeartBeat := time.Now()

		if int(latestHeartBeat.Sub(mgr.lastHeartBeat))/1000000000 > 3 {
			mgr.connectionAlive <- true
			return
		}

		err := mgr.conn.WriteMessage(websocket.TextMessage, []byte(mgr.worker_utilization()))
		if err != nil {
			log.Println("Write Error", err)
			return
		}

		msgType, bytes, err := mgr.conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket closed.")
			return
		}

		if msg := string(bytes[:]); msgType != websocket.TextMessage && msg != "ack" {
			log.Println("Unrecognized message received.")
			return
		} else {
			mgr.lastHeartBeat = time.Now()
		}

		time.Sleep(3 * time.Second)
	}
}


func publish(job *Job) {
	b, err := json.Marshal(&job)

	if err != nil {
		fmt.Println("error in pushing job")
	}

	_, err = http.Post(fmt.Sprintf("http://%s/publish", server),
		"application/json", bytes.NewBuffer(b))

	if err != nil {
		fmt.Println("error in response")
	}
}

func main() {
	mgr, err := NewManager(5)
	if err != nil {
		fmt.Println(err)
		return
	}

	job := &Job {
		Name: "Add",
		Jid: "123",
		Args: map[string]interface{}{
			"Length": 10,
			"Width": 10,
		},
		Queue: "Normal",
	}
	publish(job)
	mgr.Process()
}
