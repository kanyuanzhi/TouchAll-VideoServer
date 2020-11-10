package main

import (
	"TouchAll-VideoServer/socket"
	"TouchAll-VideoServer/websocket"
	"time"
)

type VideoServer struct {
	SocketServer *socket.SocketServer
	WsClients    *websocket.WsClients
	WsServer     *websocket.WsServer
}

func NewVideoServer() *VideoServer {
	wsClients := websocket.NewClients()
	return &VideoServer{
		SocketServer: socket.NewSocketServer(wsClients),
		WsServer:     websocket.NewWsServer(wsClients),
		WsClients:    wsClients,
	}
}

func (videoServer *VideoServer) Start() {
	go videoServer.WsClients.Start()
	//go videoServer.WsClients.Status()
	go videoServer.WsServer.Start()
	go videoServer.SocketServer.Start()
}

func main() {
	videoServer := NewVideoServer()
	videoServer.Start()
	for {
		time.Sleep(time.Second)
	}
}
