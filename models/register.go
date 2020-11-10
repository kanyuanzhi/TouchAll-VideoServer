package models

import "github.com/gorilla/websocket"

type Register struct {
	CameraID int
	Conn     *websocket.Conn
}

func NewRegister(cameraID int, conn *websocket.Conn) *Register {
	return &Register{
		CameraID: cameraID,
		Conn:     conn,
	}
}
