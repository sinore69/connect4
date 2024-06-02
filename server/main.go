package main

import (
	"log"
	"net/http"
	generator "server/generate"
	TypeOf "server/types"

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
	//	go reader(conn)
	conn.WriteJSON(roomId)
	log.Println(room.ActiveRooms)
}
func reader(conn *websocket.Conn, session *TypeOf.Players) {
	for {
		var board TypeOf.Board
		err := conn.ReadJSON(&board)
		if err != nil {
			panic(err)
		}
		log.Println(board)
		//do board logic
		// combine both write logic into oneboard
		// and call inside reader function
		newBoard := updateState(&board)
		writer(session, newBoard)
	}
}
func updateState(board *TypeOf.Board) *TypeOf.Board {
	board.Board[board.LastMove.RowIndex][board.LastMove.ColIndex] = 1
	return board
}
func writer(session *TypeOf.Players, newBoard *TypeOf.Board) {
	creator, player := session.Creator, session.Player
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
