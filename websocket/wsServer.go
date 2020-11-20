package websocket

import (
	"TouchAll-VideoServer/models"
	"TouchAll-VideoServer/utils"
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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
	config := utils.NewConfig()
	port := config.GetWebSocketConfig().(string)
	addr := flag.String("addr", ":"+port, "http service address")
	http.HandleFunc("/ws", wsServer.serveWs)

	log.Printf("Start WsServer the of video server on port %s", port)

	_ = http.ListenAndServe(*addr, nil)
}

func (wsServer *WsServer) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			_ = conn.Close()
			return
		}
		go wsServer.handleConn(message, conn)
	}
}
func (wsServer *WsServer) handleConn(message []byte, conn *websocket.Conn) {
	var request models.Request
	err := json.Unmarshal(message, &request)
	if err != nil {
		log.Println(err.Error())
		_ = conn.Close()
		return
	}

	switch request.RequestType {
	case 50:
		previousCameraID := request.PreviousCameraID
		if previousCameraID == -1 {
			// 前端无切换请求不同监控的需求，用于数据墙
			register := models.NewRegister(request.CameraID, conn)
			wsServer.wsClients.individualRegister <- register
		} else {
			// 满足前端通过点击切换请求不同监控
			if previousCameraID != 0 {
				delete(wsServer.wsClients.members[previousCameraID], conn)
				if len(wsServer.wsClients.members[previousCameraID]) == 0 {
					delete(wsServer.wsClients.members, previousCameraID)
				}
			} else if previousCameraID == request.CameraID {
				return
			}
			register := models.NewRegister(request.CameraID, conn)
			wsServer.wsClients.register <- register
		}
	case 53:
		previousCameraID := request.PreviousCameraID
		if previousCameraID == -1 {
			aiRegister := models.NewRegister(request.CameraID, conn)
			wsServer.wsClients.individualAIRegister <- aiRegister
		} else {
			if previousCameraID := request.PreviousCameraID; previousCameraID != 0 {
				delete(wsServer.wsClients.aiMembers[previousCameraID], conn)
				if len(wsServer.wsClients.aiMembers[previousCameraID]) == 0 {
					delete(wsServer.wsClients.aiMembers, previousCameraID)
				}
			} else if previousCameraID == request.CameraID {
				return
			}
			aiRegister := models.NewRegister(request.CameraID, conn)
			wsServer.wsClients.aiRegister <- aiRegister
		}
	default:
		return
	}
}
