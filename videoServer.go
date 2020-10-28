package main

import "time"

type VideoServer struct {
	SocketServer *SocketServer
	WsClients    *WsClients
	WsServer     *WsServer
}

func NewVideoServer() *VideoServer {
	wsClients := NewClients()
	return &VideoServer{
		SocketServer: NewSocketServer(wsClients),
		WsServer:     NewWsServer(wsClients),
		WsClients:    wsClients,
	}
}

func (vs *VideoServer) Start() {
	go vs.WsClients.Start()
	go vs.WsServer.Start()
	go vs.SocketServer.Start()
	for {
		time.Sleep(time.Second)
	}
}
