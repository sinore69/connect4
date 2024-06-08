package types

import (
	"github.com/gorilla/websocket"
)

type IncorrectRoomIid struct {
	Message string
}
type RoomId struct {
	Id int
}

type Players struct {
	Creator        *websocket.Conn
	DisableCreator bool
	Player         *websocket.Conn
	DisablePlayer  bool
	DummyReader    chan bool
}
type LastMove struct {
	RowIndex int
	ColIndex int
}
