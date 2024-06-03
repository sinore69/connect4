package main

import (
	"log"
	"net/http"
	generator "server/generate"
	TypeOf "server/types"
	//Game "server/game"
	"github.com/gorilla/websocket"
)

type Rooms struct {
	ActiveRooms map[int]TypeOf.Players
}
type Board struct {
	Board     [10][10]int
	MoveCount int
	LastMove  TypeOf.LastMove
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
func (board *Board) reader(conn *websocket.Conn, session *TypeOf.Players) {
	for {
		err := conn.ReadJSON(&board)
		if err != nil {
			panic(err)
		}
		newBoard := UpdateState(board)
		writer(session, newBoard)
	}
}
func UpdateState(board *Board) *Board {
	board.MoveCount++
	if board.MoveCount%2 == 0 {
		board.Board[board.LastMove.RowIndex][board.LastMove.ColIndex] = 1
	} else {
		board.Board[board.LastMove.RowIndex][board.LastMove.ColIndex] = 2
	}
	return board
}
func writer(session *TypeOf.Players, newBoard *Board) {
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
	board:=Board{}
	go board.reader(session.Creator, &session)
	go board.reader(session.Player, &session)
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
