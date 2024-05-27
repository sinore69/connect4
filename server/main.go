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
type Rooms struct {
	ActiveRooms map[int]players
}
type players struct {
	Creator chan board
	player  chan board
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
	defer conn.Close()
	room.ActiveRooms[generator.NewRoomId()] = players{
		Creator: make(chan board),
	}
	log.Println(room.ActiveRooms)
	for {
		var msg App
		err := conn.ReadJSON(&msg)
		if err != nil {
			panic(err)
		}
		log.Println(msg)
	}
}
func (room *Rooms) JoinRoom(w http.ResponseWriter, r *http.Request) {
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}

func main() {
	var room Rooms = *NewRoom()
	http.HandleFunc("/echo", room.CreateRoom)
	http.HandleFunc("/echo", room.JoinRoom)
	log.Println("Server started at :5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
