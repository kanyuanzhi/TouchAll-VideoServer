package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"videoServer/models"
)

const MAXMESSAGESIZE = 4024

var upgrader = websocket.Upgrader{
	ReadBufferSize:  MAXMESSAGESIZE,
	WriteBufferSize: MAXMESSAGESIZE,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket服务器
type WsServer struct {
	wsClients *WsClients
}

func NewWsServer(wsClients *WsClients) *WsServer {
	return &WsServer{
		wsClients: wsClients,
	}
}

func (wsServer *WsServer) Start() {
	addr := flag.String("addr", "localhost:9081", "http service address")
	http.HandleFunc("/ws", wsServer.serveWs)

	log.Println("start wsServer on port 9081")

	http.ListenAndServe(*addr, nil)
}

func (wsServer *WsServer) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		conn.Close()
		return
	}

	var request models.Request
	err = json.Unmarshal(message, &request)
	if err != nil {
		log.Println("json:", err)
		conn.Close()
		return
	}

	register := models.NewRegister(request.Camera, conn)
	wsServer.wsClients.register <- register
}
