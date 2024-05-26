package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)
type App struct{
	Message string
	Num int
}
type Rooms struct{
	Conn []int
}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func (room *Rooms)echoHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
		panic(err)
    }
	i:=len(room.Conn)
	room.Conn=append(room.Conn,i)
	log.Println(room.Conn)
	defer conn.Close()
	for{
		var msg App
		err:=conn.ReadJSON(&msg);
		if err!=nil{
			panic(err)
		}
		log.Println(msg)
	}
}


func main() {
	var room Rooms
	http.HandleFunc("/echo", room.echoHandler)
	log.Println("Server started at :5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}