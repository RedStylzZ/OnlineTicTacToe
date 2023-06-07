package ttt

import (
	"fmt"
	"strings"
)

const emptyField = "-"

type TicTacToe struct {
	field         [3][3]string
	currentPlayer *Player
}

func NewTicTacToe(first, second string) *TicTacToe {
	ttt := &TicTacToe{}
	player1 := &Player{
		id:   1,
		name: first,
		char: "X",
	}
	player2 := &Player{
		id:         2,
		name:       second,
		char:       "O",
		nextPlayer: player1,
	}
	player1.nextPlayer = player2
	// ttt.players = players
	ttt.currentPlayer = player1
	ttt.initializeField()
	return ttt
}

func (t *TicTacToe) Start() {
	t.PrintField()
}

func (t *TicTacToe) initializeField() {
	for x := 0; x < len(t.field); x++ {
		for y := 0; y < len(t.field[x]); y++ {
			t.field[y][x] = emptyField
		}
	}
}

func (t *TicTacToe) PrintField() {
	fmt.Println("   0 1 2")
	for i := range t.field {
		row := strings.Join(t.field[i][:], "|")
		fmt.Printf("%d |%s|\n", i, row)
	}
}

func (t *TicTacToe) SetField(x, y int) error {
	if t.field[x][y] != emptyField {
		return fmt.Errorf("field [%d][%d] is not empty", x, y)
	}
	t.field[x][y] = t.currentPlayer.GetChar()
	return nil
}

func (t *TicTacToe) SwitchPlayers() {
	t.currentPlayer = t.currentPlayer.nextPlayer
}

func (t *TicTacToe) GetCurrentPlayer() *Player {
	return t.currentPlayer
}

func (t *TicTacToe) CheckEnd() bool {
	char := t.currentPlayer.GetChar()

	var rowCount, columnCount int
	for row := range t.field {
		rowCount, columnCount = 0, 0
		for column := range t.field[row] {
			if t.field[row][column] == char {
				if rowCount++; rowCount == 3 {
					return true
				}
			}
			if t.field[column][row] == char {
				if columnCount++; columnCount == 3 {
					return true
				}
			}
		}
	}

	var leftCount, rightCount int
	for x := 0; x < len(t.field); x++ {
		if t.field[x][x] == char {
			if leftCount++; leftCount == 3 {
				return true
			}
		}
		if t.field[x][len(t.field)-1-x] == char {
			if rightCount++; rightCount == 3 {
				return true
			}
		}
	}

	return false
}
