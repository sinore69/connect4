package game

import (
	TypeOf "server/types"
)

func CheckWin(board TypeOf.Board) bool {
	//row:=board.LastMove.RowIndex
	//col:=board.LastMove.ColIndex
	//horizontal check
	//vertical check
	//diagonal check
	return false
}

func UpdateState(board *TypeOf.Board) *TypeOf.Board {
	board.MoveCount++
	if board.MoveCount%2 == 0 {
		board.Board[board.LastMove.RowIndex][board.LastMove.ColIndex] = 1
	} else {
		board.Board[board.LastMove.RowIndex][board.LastMove.ColIndex] = 2
	}
	return board
}
