package game

func BoardCompleted(board [10][10]int)bool{
	for i:=0;i<10;i++{
		for j:=0;j<10;j++{
			if board[i][j]==0{
				return false
			}
		}
	}
	return true
}

func Checkwin(board [10][10]int,row int,col int) bool {
	if horizontalwin(board, row, col) {
		return true
	}
	if verticalwin(board, row, col) {
		return true
	}
	if diagonalwin(board, row, col) {
		return true
	}
	return false
}

func diagonalwin(board [10][10]int, row int, col int) bool {
	leftcount,rightcount:=1,1
	leftcount+=countDiscs(board,row,col,1,1)
	leftcount+=countDiscs(board,row,col,-1,-1)
	rightcount+=countDiscs(board,row,col,1,-1)
	rightcount+=countDiscs(board,row,col,-1,1)
	return rightcount>=4||leftcount>=4
}

func verticalwin(board [10][10]int, row int, col int) bool {
	count:=1
	count+=countDiscs(board,row,col,1,0)
	count+=countDiscs(board,row,col,-1,0)
	return count>=4
}

func horizontalwin(board [10][10]int, row int, col int) bool {
	count:=1
	count+=countDiscs(board,row,col,0,1)
	count+=countDiscs(board,row,col,0,-1)
	return count>=4
}

func countDiscs(board [10][10]int,row, column, dRow, dCol int) int {
    player := board[row][column]
    count := 0
    for {
        row += dRow
        column += dCol
        if player==0|| row < 0 || row >= len(board) || column < 0 || column >= len(board[0]) || board[row][column] != player {
            break
        }
        count++
    }
    return count
}