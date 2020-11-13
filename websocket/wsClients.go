package websocket

import (
	"TouchAll-VideoServer/models"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type WsClients struct {
	members    map[int]map[*websocket.Conn]bool
	aiMembers  map[int]map[*websocket.Conn]bool
	register   chan *models.Register
	aiRegister chan *models.Register
	Video      chan *models.Video
	AIVideo    chan *models.Video
}

func NewClients() *WsClients {
	return &WsClients{
		members:    make(map[int]map[*websocket.Conn]bool),
		aiMembers:  make(map[int]map[*websocket.Conn]bool),
		register:   make(chan *models.Register),
		aiRegister: make(chan *models.Register),
		Video:      make(chan *models.Video),
		AIVideo:    make(chan *models.Video),
	}
}

func (wsClients *WsClients) Start() {
	for {
		select {
		case request := <-wsClients.register: // 注册对普通摄像机的请求
			if _, has := wsClients.members[request.CameraID]; !has {
				wsClients.members[request.CameraID] = make(map[*websocket.Conn]bool)
			}
			wsClients.members[request.CameraID][request.Conn] = true
		case request := <-wsClients.aiRegister: // 注册对AI摄像机的请求
			if _, has := wsClients.aiMembers[request.CameraID]; !has {
				wsClients.aiMembers[request.CameraID] = make(map[*websocket.Conn]bool)
			}
			wsClients.aiMembers[request.CameraID][request.Conn] = true

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
		case video := <-wsClients.AIVideo:
			if aiMembers, has := wsClients.aiMembers[video.CameraID]; has {
				for aiMember := range aiMembers {
					go func(aiMember *websocket.Conn) {
						err := aiMember.WriteMessage(websocket.BinaryMessage, video.Image)
						if err != nil {
							log.Printf("write errro: %s", err.Error())
							aiMember.Close()
							delete(aiMembers, aiMember)
							if len(aiMembers) == 0 {
								delete(wsClients.aiMembers, video.CameraID)
							}
						}
					}(aiMember)
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
