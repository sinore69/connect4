package main

import (
	"log"
	"net/http"
	generator "server/generate"
	TypeOf "server/types"
	Game "server/game"
	"github.com/gorilla/websocket"
)

type Rooms struct {
	ActiveRooms map[int]TypeOf.Players
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewRoom() *Rooms {
	return &Rooms{
		ActiveRooms: make(map[int]TypeOf.Players),
	}
}
func (room *Rooms) CreateRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	id := generator.NewRoomId()
	room.ActiveRooms[id] = TypeOf.Players{
		Creator: conn,
	}
	roomId := &TypeOf.RoomId{
		Id: id,
	}
	conn.WriteJSON(roomId)
	log.Println(room.ActiveRooms)
}
func reader(conn *websocket.Conn, session *TypeOf.Players) {
	var board TypeOf.Board
	for {
		err := conn.ReadJSON(&board)
		if err != nil {
			panic(err)
		}
		newBoard := Game.UpdateState(&board)
		writer(session, newBoard)
	}
}

func writer(session *TypeOf.Players, newBoard *TypeOf.Board) {
	creator, player := session.Creator, session.Player
	log.Println(newBoard)
	creator.WriteJSON(&newBoard)
	player.WriteJSON(&newBoard)
}
func (room *Rooms) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	var id TypeOf.RoomId
	err = conn.ReadJSON(&id)
	if err != nil {
		panic(err)
	}
	session := room.ActiveRooms[id.Id]
	if session.Creator == nil {
		error := &TypeOf.IncorrectRoomIid{
			Message: "No Such Room Id",
		}
		conn.WriteJSON(error)
		return
	} else {
		conn.WriteJSON(id)
	}
	session.Player = conn
	room.ActiveRooms[id.Id] = session
	log.Println(room.ActiveRooms)
	go reader(session.Creator, &session)
	go reader(session.Player, &session)
}
func main() {
	var room Rooms = *NewRoom()
	http.HandleFunc("/create", room.CreateRoom)
	http.HandleFunc("/join", room.JoinRoom)
	log.Println("Server started at :5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
