package websocket

import (
	"TouchAll-VideoServer/models"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type WsClients struct {
	members map[int]map[*websocket.Conn]bool
	//lockers  map[*websocket.Conn]*sync.Mutex
	register   chan *models.Register
	unRegister chan *models.Register
	Video      chan *models.Video
}

func NewClients() *WsClients {
	return &WsClients{
		//lockers:  make(map[*websocket.Conn]*sync.Mutex),
		members:  make(map[int]map[*websocket.Conn]bool),
		register: make(chan *models.Register),
		Video:    make(chan *models.Video),
	}
}

func (wsClients *WsClients) Start() {
	for {
		select {
		case request := <-wsClients.register: // 注册
			if _, has := wsClients.members[request.CameraID]; !has {
				wsClients.members[request.CameraID] = make(map[*websocket.Conn]bool)
			}
			wsClients.members[request.CameraID][request.Conn] = true

		case video := <-wsClients.Video:
			if members, has := wsClients.members[video.CameraID]; has {
				for member := range members {
					go func(member *websocket.Conn) {
						err := member.WriteMessage(websocket.BinaryMessage, video.Image)
						if err != nil {
							log.Printf("write errro: %s", err.Error())
							member.Close()
							delete(members, member)
							if len(members) == 0 {
								delete(wsClients.members, video.CameraID)
							}
						}
					}(member)
				}
			}
		}
	}
}

func (wsClients *WsClients) Status() {
	for {
		log.Println(wsClients.members)
		time.Sleep(2 * time.Second)
	}
}
