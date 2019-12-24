package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"time"
	"fmt"
	"encoding/json"
	"bytes"
	"log"
)

var server = "localhost:3000"

type Manager struct {
	workers int
	conn *websocket.Conn
	connectionAlive chan bool
	lastHeartBeat time.Time
}


type Job struct {
	Jid string
	Name string
	Args []interface{}
}

func NewManager(concurrency int) *Manager {
	return &Manager {
		workers: concurrency,
		connectionAlive: make(chan bool),
		lastHeartBeat: time.Now(),
	}
}

func (mgr *Manager) Process() {

	uri := "ws://localhost:3000/cosume"
	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)

	fmt.Println(conn)

	if err != nil {
		fmt.Println(err)
		return
	}

	mgr.conn = conn

	go mgr.heartbeat_server()

	select {
	case <-mgr.connectionAlive:
		fmt.Println("worker is dead")
	}
}


func (mgr *Manager) heartbeat_server() {
	defer mgr.conn.Close()

	latestHeartBeat := time.Now()
	if int(latestHeartBeat.Sub(mgr.lastHeartBeat))/1000000000 > 60 {
		mgr.connectionAlive <- true
		return
	}

	err := mgr.conn.WriteMessage(websocket.TextMessage, []byte("utilization-50"))
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
