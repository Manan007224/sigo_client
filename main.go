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
)

var server = "localhost:3000"

type Manager struct {
	pool []*Worker
	workers int
	conn *websocket.Conn
	connectionAlive chan bool
	lastHeartBeat time.Time
	busyWorkers *Set
	freeWorkers *Set
}

type Worker struct {
	jobChan chan Job
}


type Job struct {
	Jid string
	Name string
	Args []interface{}
}

func NewManager(concurrency int) *Manager {
	mgr:= &Manager {
		pool: []*Worker{},
		workers: concurrency,
		connectionAlive: make(chan bool),
		lastHeartBeat: time.Now(),
	}

	uri := "ws://localhost:3000/consume"
	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	mgr.conn = conn

	for w := 0; w < mgr.workers; w++ {
		worker := &Worker {
			jobChan: make(chan Job)
		}
		mgr.pool = append(mgr.pool, worker)
		freeWorkers.Add(worker)
	}

}

func (mgr *Manager) Process() {
	go mgr.heartbeat_server()

	select {
	case <-mgr.connectionAlive:
		fmt.Println("worker is dead")
	}
}

func (mgr *Manager) worker_utilization() {
	fmt.Sprintf("utilization-%s", strconv.Itoa(mgr.busyWorkers) / len(mgr.pool))
}


func (mgr *Manager) heartbeat_server() {
	for {
		latestHeartBeat := time.Now()

		if int(latestHeartBeat.Sub(mgr.lastHeartBeat))/1000000000 > 3 {
			mgr.connectionAlive <- true
			return
		}

		err := mgr.conn.WriteMessage(websocket.TextMessage, []byte(worker_utilization()))
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
			log.Println("Received: ack.")
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
	mgr := NewManager(5)
	mgr.Process()
}
