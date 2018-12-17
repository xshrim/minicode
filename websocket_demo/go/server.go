package main

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

var gws *websocket.Conn

func Echo(ws *websocket.Conn) {
	gws = ws
	for {
		var req string
		if err := websocket.Message.Receive(ws, &req); err != nil {
			log.Println("receive error:", err)
		} else {
			log.Println(req)
			Deal(req, ws)
		}
	}
}

func Deal(req string, ws *websocket.Conn) {
	switch req {
	case "hello":
		if err := websocket.Message.Send(ws, "welcome"); err != nil {
			log.Println("send error:", err)
			break
		}
	default:
		if err := websocket.Message.Send(ws, req); err != nil {
			log.Println("send error:", err)
		}
	}
}

func Tick() {
	for {
		time.Sleep(5 * time.Second)
		if gws != nil {
			if err := websocket.Message.Send(gws, "Tick"); err != nil {
				log.Println("send error:", err)
			}
		}
	}
}

func test(x int) {
	log.Println(x)
}

type Handler func(int)

func main() {

	ws := websocket.Handler(Echo)
	http.Handle("/ws", ws)
	//http.Handle("/ws", websocket.Handler(Echo))
	go Tick()
	http.ListenAndServe("127.0.0.1:1234", nil)
}
