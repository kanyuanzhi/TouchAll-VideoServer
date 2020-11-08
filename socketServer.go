package main

import (
	"TouchAll-VideoServer/models"
	"TouchAll-VideoServer/protocal"
	"TouchAll-VideoServer/utils"
	"fmt"
	"log"
	"net"
)

// 接口服务器
type SocketServer struct {
	wsClients *WsClients
}

func NewSocketServer(wsClients *WsClients) *SocketServer {
	return &SocketServer{
		wsClients: wsClients,
	}
}

func (socketServer *SocketServer) Start() {
	config := utils.NewConfig()
	port := config.GetSocketConfig().(string)
	l, err := net.Listen("tcp", ":"+port)
	log.Printf("Start the SocketServer of video server on port %s", port)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}
		go socketServer.handleConn(conn)
	}
}

// 处理socket连接，完成数据解包
func (socketServer *SocketServer) handleConn(conn net.Conn) {
	defer conn.Close()
	tempBuffer := make([]byte, 0)
	readerChannel := make(chan []byte, 1024)
	go socketServer.reader(readerChannel)
	for {
		var buffer = make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		tempBuffer = protocal.Unpack(append(tempBuffer, buffer[:n]...), readerChannel)
	}
}

func (socketServer *SocketServer) reader(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			video := new(models.Video)
			video.Camera = protocal.BytesToInt(data[:4])
			video.Image = data[4:]
			socketServer.wsClients.video <- video
		}
	}
}
