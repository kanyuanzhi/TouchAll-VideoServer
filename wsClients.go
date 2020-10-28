package main

import (
	"github.com/gorilla/websocket"
	"log"
	"videoServer/models"
)

type WsClients struct {
	members map[int]map[*websocket.Conn]bool
	//lockers  map[*websocket.Conn]*sync.Mutex
	register chan *models.Register
	video    chan *models.Video
}

func NewClients() *WsClients {
	return &WsClients{
		//lockers:  make(map[*websocket.Conn]*sync.Mutex),
		members:  make(map[int]map[*websocket.Conn]bool),
		register: make(chan *models.Register),
		video:    make(chan *models.Video),
	}
}

func (wsClients *WsClients) Start() {
	for {
		select {
		case request := <-wsClients.register: // 注册
			if _, has := wsClients.members[request.Camera]; !has {
				wsClients.members[request.Camera] = make(map[*websocket.Conn]bool)
			}
			wsClients.members[request.Camera][request.Conn] = true
		case video := <-wsClients.video:
			if members, has := wsClients.members[video.Camera]; has {
				for member := range members {
					go func(member *websocket.Conn) {
						err := member.WriteMessage(websocket.BinaryMessage, video.Image)
						if err != nil {
							log.Printf("write errro: %s,%s", err)
							member.Close()
							delete(members, member)
						}
					}(member)
				}
			}
		}
	}
}
