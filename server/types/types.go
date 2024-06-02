package types

import "github.com/gorilla/websocket"

type IncorrectRoomIid struct {
	Message string
}
type RoomId struct {
	Id int
}

type Players struct {
	Creator *websocket.Conn
	Player  *websocket.Conn
}
type Board struct {
	Board     [10][10]int
	MoveCount int
	LastMove  LastMove
}
type LastMove struct {
	RowIndex int
	ColIndex int
}
