package main

import (
	"log"
	"net/http"
	"server/game"
	generator "server/generate"
	TypeOf "server/types"

	"github.com/gorilla/websocket"
)

type Message struct {
	Text string
}

type InitialState struct {
	Disable bool
}

type Rooms struct {
	ActiveRooms map[int]TypeOf.Players
}

type Board struct {
	Board     [10][10]int
	MoveCount int
	Disable   bool
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
		newBoard := UpdateState(board, session)
		if game.Checkwin(newBoard.Board, newBoard.LastMove.RowIndex, newBoard.LastMove.ColIndex) {
			endgame(session, newBoard, false)
			break
		}
		if game.BoardCompleted(newBoard.Board) {
			log.Println("error")
			endgame(session, newBoard, true)
			break
		}
		writer(session, newBoard)
	}
}

func endgame(session *TypeOf.Players, newBoard *Board, boardcompleted bool) {
	freezestate(session)
	writer(session, newBoard)
	if boardcompleted {
		drawmsg := Message{
			Text: "Draw",
		}
		session.Creator.WriteJSON(drawmsg)
		session.Player.WriteJSON(drawmsg)
		return
	}
	winningmsg := Message{
		Text: "You Won",
	}
	losingmsg := Message{
		Text: "Your Opponent Won",
	}
	if newBoard.MoveCount%2 == 0 {
		session.Creator.WriteJSON(losingmsg)
		session.Player.WriteJSON(winningmsg)
	} else {
		session.Creator.WriteJSON(winningmsg)
		session.Player.WriteJSON(losingmsg)
	}
}

func freezestate(session *TypeOf.Players) {
session.DisableCreator = true
session.DisablePlayer = true
}

func UpdateState(board *Board, session *TypeOf.Players) *Board {
	if board.Board[board.LastMove.RowIndex][board.LastMove.ColIndex] != 0 {
		return board
	}
	board.MoveCount++
	if board.MoveCount%2 == 0 {
		board.Board[board.LastMove.RowIndex][board.LastMove.ColIndex] = 1
		session.DisableCreator = false
		session.DisablePlayer = true
	} else {
		board.Board[board.LastMove.RowIndex][board.LastMove.ColIndex] = 2
		session.DisableCreator = true
		session.DisablePlayer = false
	}
	return board
}

func writer(session *TypeOf.Players, newBoard *Board) {
	creator, player := session.Creator, session.Player
	log.Println(newBoard)
	log.Println(session)
	newBoard.Disable = session.DisableCreator
	creator.WriteJSON(&newBoard)
	newBoard.Disable = session.DisablePlayer
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
		conn.WriteJSON(&id)
	}
	session.Player = conn
	room.ActiveRooms[id.Id] = session
	log.Println(room.ActiveRooms)
	board := Board{}
	initialstate(&session)
	go board.reader(session.Creator, &session)
	go board.reader(session.Player, &session)
}

func initialstate(session *TypeOf.Players) {
	disable := InitialState{
		Disable: false,
	}
	session.Creator.WriteJSON(disable)
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
