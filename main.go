package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"time"
	"fmt"
	"encoding/json"
	"bytes"
	"log"
	uuid "github.com/google/uuid"
)

var server = "localhost:3000"

type Manager struct {
	workers map[uuid.UUID]*Worker
	concurrency int
	aliveWorkers uint32
	deadWorkerChan chan *Worker
}

type Worker struct {
	conn *websocket.Conn
	pid uuid.UUID
	lastHeartBeat time.Time
	deadWorkerChan chan *Worker
}

type Job struct {
	Jid string
	Name string
	Args []interface{}
}

func NewManager(concurrency int) *Manager {
	return &Manager {
		workers: make(map[uuid.UUID]*Worker),
		concurrency: concurrency,
		deadWorkerChan: make(chan *Worker),
	}
}

func (mgr *Manager) Process() {
	for w := 0; w < mgr.concurrency; w++ {
		worker := &Worker {
			pid: uuid.New(),
			deadWorkerChan: mgr.deadWorkerChan,
		}

		uri := fmt.Sprintf("ws://%s/cosume?pid=%s", server, worker.pid)
		conn, _, err := websocket.DefaultDialer.Dial(uri, make(http.Header))

		if err != nil {
			fmt.Println("error in establishing websocket connection to the server")
			fmt.Println(err)
			continue
		}
		worker.conn = conn
		mgr.workers[worker.pid] = worker
		worker.Run()
	}
}

func (w *Worker) Run() {
	// continously pulling jobs from the master
	go w.heartbeat_server()

	// select {
	// case <-deadWorkerChan:
	// 	// handle client_worker_close connection here
	// }
}

func (w *Worker) heartbeat_server() {
	defer w.conn.Close()

	latestHeartBeat := time.Now()
	if int(latestHeartBeat.Sub(w.lastHeartBeat))/1000000000 > 60 {
		w.deadWorkerChan <- w
		return
	}

	err := w.conn.WriteMessage(websocket.TextMessage, []byte("ping"))
	if err != nil {
		log.Println("Write Error", err)
		return
	}

	msgType, bytes, err := w.conn.ReadMessage()
	if err != nil {
		log.Println("WebSocket closed.")
		return
	}

	if msg := string(bytes[:]); msgType != websocket.TextMessage && msg != "pong" {
		log.Println("Unrecognized message received.")
		return
	} else {
		w.lastHeartBeat = time.Now()
		log.Println("Received: pong.")
	}

	time.Sleep(60 * time.Second)

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
