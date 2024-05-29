package main

import (
	"log"
	"net/http"
	generator "server/generate"

	"github.com/gorilla/websocket"
)

type IncorrectRoomIid struct {
	Message string
}
type RoomId struct {
	Id int
}
type Rooms struct {
	ActiveRooms map[int]players
}
type players struct {
	Creator chan board
	Player  chan board
}
type board struct {
	board [][]int
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewRoom() *Rooms {
	return &Rooms{
		ActiveRooms: make(map[int]players),
	}
}
func (room *Rooms) CreateRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	ch := make(chan board)
	id := generator.NewRoomId()
	room.ActiveRooms[id] = players{
		Creator: ch,
	}
	roomId := &RoomId{
		Id: id,
	}
	readerCh := (chan<- board)(ch)
	writerCh := (<-chan board)(ch)
	go reader(&readerCh, conn)
	go writer(&writerCh, conn)
	conn.WriteJSON(roomId)
	log.Println(room.ActiveRooms)
}
func reader(ch *chan<- board, conn *websocket.Conn) {
	for {
		var msg IncorrectRoomIid
		err := conn.ReadJSON(&msg)
		if err != nil {
			panic(err)
		}
		log.Println(msg)
	}
}
func writer(ch *<-chan board, conn *websocket.Conn) {

}
func (room *Rooms) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	var id RoomId
	err = conn.ReadJSON(&id)
	if err != nil {
		panic(err)
	}
	session:= room.ActiveRooms[id.Id]
	if session.Creator == nil {
		error:=&IncorrectRoomIid{
			Message: "No Such Room Id",
		}
		conn.WriteJSON(error)
		return
	}else{
		conn.WriteJSON(id)
	}
	ch := make(chan board)
	session.Player = ch
	room.ActiveRooms[id.Id] = session
	log.Println(room.ActiveRooms)
	readerCh := (chan<- board)(ch)
	writerCh := (<-chan board)(ch)
	go reader(&readerCh, conn)
	go writer(&writerCh, conn)
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
