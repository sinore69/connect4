package main

import (
	"context"
	"log"
	"net/http"
	"server/game"
	generator "server/generate"
	TypeOf "server/types"
	"sync"
	"time"

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
	Id        int
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
var mutex sync.Mutex

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
	done := make(chan bool)
	mutex.Lock()
	room.ActiveRooms[id] = TypeOf.Players{
		Creator:     conn,
		DummyReader: done,
	}
	mutex.Unlock()
	roomId := &TypeOf.RoomId{
		Id: id,
	}
	conn.WriteJSON(roomId)
	go dummyreader(done, conn, room, id)
	log.Println(room.ActiveRooms)
}

func dummyreader(done <-chan bool, conn *websocket.Conn, room *Rooms, id int) {
	start := time.Now()
	var elapsed time.Duration
outer:
	for {
		elapsed = time.Since(start)
		select {
		case <-done:
			log.Println("done")
			break outer
		default:
			if elapsed > 3*time.Minute {
				log.Println("connection timeout")
				conn.Close()
				mutex.Lock()
				delete(room.ActiveRooms, id)
				mutex.Unlock()
				break outer
			}
		}
	}
}

func (board *Board) reader(conn *websocket.Conn, room *Rooms, done chan<- bool, ctx context.Context) {
	defer conn.Close()
	defer mutex.Unlock()
	defer delete(room.ActiveRooms, board.Id)
	defer mutex.Lock()
outer:
	for {
		select {
		case <-ctx.Done():
			break outer
		default:
			err := conn.ReadJSON(&board)
			if err != nil {
				log.Println("connection lost")
				done <- true
				break outer
			}
			session := room.ActiveRooms[board.Id]
			newBoard := UpdateState(board, &session)
			if game.Checkwin(newBoard.Board, newBoard.LastMove.RowIndex, newBoard.LastMove.ColIndex) {
				endgame(&session, newBoard, false)
				break outer
			}
			if game.BoardCompleted(newBoard.Board) {
				log.Println("board is complet")
				endgame(&session, newBoard, true)
				break outer
			}
			writer(&session, newBoard)
		}
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
	mutex.Lock()
	session := room.ActiveRooms[id.Id]
	mutex.Unlock()
	if session.Creator == nil {
		error := &TypeOf.IncorrectRoomIid{
			Message: "No Such Room Id",
		}
		conn.WriteJSON(error)
		return
	} else {
		conn.WriteJSON(&id)
	}
	session.DummyReader <- true
	close(session.DummyReader)
	session.Player = conn
	mutex.Lock()
	room.ActiveRooms[id.Id] = session
	mutex.Unlock()
	log.Println(room.ActiveRooms)
	board := Board{
		Id: id.Id,
	}
	initialstate(&session)
	done := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	go board.reader(session.Creator, room, done, ctx)
	go board.reader(session.Player, room, done, ctx)
	<-done
	cancel()
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
