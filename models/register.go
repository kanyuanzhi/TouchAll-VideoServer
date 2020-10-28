package models

import "github.com/gorilla/websocket"

type Register struct {
	Camera int
	Conn   *websocket.Conn
}

func NewRegister(camera int, conn *websocket.Conn) *Register {
	return &Register{
		Camera: camera,
		Conn:   conn,
	}
}
