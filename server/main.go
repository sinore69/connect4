package main

import (
	"log"
	"net/http"
	generator "server/generate"

	"github.com/gorilla/websocket"
)

type App struct {
	Message string
	Num     int
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
	ch:=make(chan board)
	room.ActiveRooms[generator.NewRoomId()] = players{
		Creator: ch,
	}
	readerCh := (chan<- board)(ch)
	go reader(&readerCh,conn)
	log.Println(room.ActiveRooms)
}
func reader(ch *chan<- board,conn *websocket.Conn){
	for {
		var msg App
		err := conn.ReadJSON(&msg)
		if err != nil {
			panic(err)
		}
		log.Println(msg)
	}
}
func writer(){

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
	 	session:=room.ActiveRooms[id.Id]
		if session.Creator==nil{
			panic("no creator")
		}
		ch:=make(chan board)
		session.Player=ch
		room.ActiveRooms[id.Id]=session
		log.Println(room.ActiveRooms)
		readerCh:=(chan<-board)(ch)
		go reader(&readerCh,conn)
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
