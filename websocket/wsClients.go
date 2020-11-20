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

	individualMembers    map[int]map[*websocket.Conn]bool
	individualAIMembers  map[int]map[*websocket.Conn]bool
	individualRegister   chan *models.Register
	individualAIRegister chan *models.Register

	Video   chan *models.Video
	AIVideo chan *models.Video
}

func NewClients() *WsClients {
	return &WsClients{
		members:    make(map[int]map[*websocket.Conn]bool),
		aiMembers:  make(map[int]map[*websocket.Conn]bool),
		register:   make(chan *models.Register),
		aiRegister: make(chan *models.Register),

		individualMembers:    make(map[int]map[*websocket.Conn]bool),
		individualAIMembers:  make(map[int]map[*websocket.Conn]bool),
		individualRegister:   make(chan *models.Register),
		individualAIRegister: make(chan *models.Register),

		Video:   make(chan *models.Video),
		AIVideo: make(chan *models.Video),
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
		case request := <-wsClients.individualRegister: // 注册对单个普通摄像机的请求
			if _, has := wsClients.individualMembers[request.CameraID]; !has {
				wsClients.individualMembers[request.CameraID] = make(map[*websocket.Conn]bool)
			}
			wsClients.individualMembers[request.CameraID][request.Conn] = true
		case request := <-wsClients.individualAIRegister: // 注册对单个AI摄像机的请求
			if _, has := wsClients.individualAIMembers[request.CameraID]; !has {
				wsClients.individualAIMembers[request.CameraID] = make(map[*websocket.Conn]bool)
			}
			wsClients.individualAIMembers[request.CameraID][request.Conn] = true

		case video := <-wsClients.Video:
			go wsClients.pushVideo(video, wsClients.members)
			go wsClients.pushVideo(video, wsClients.individualMembers)
		case video := <-wsClients.AIVideo:
			go wsClients.pushVideo(video, wsClients.aiMembers)
			go wsClients.pushVideo(video, wsClients.individualAIMembers)
		}
	}
}

func (wsClients *WsClients) pushVideo(video *models.Video, Members map[int]map[*websocket.Conn]bool) {
	if members, has := Members[video.CameraID]; has {
		for member := range members {
			go func(member *websocket.Conn) {
				defer func() {
					if err := recover(); err != nil {
						log.Println(err)
					}
				}()
				err := member.WriteMessage(websocket.BinaryMessage, video.Image)
				if err != nil {
					log.Printf("write errro: %s", err.Error())
					member.Close()
					delete(members, member)
					if len(members) == 0 {
						delete(wsClients.aiMembers, video.CameraID)
					}
				}
			}(member)
		}
	}
}

func (wsClients *WsClients) Status() {
	for {
		log.Println(wsClients.members)
		time.Sleep(2 * time.Second)
	}
}
