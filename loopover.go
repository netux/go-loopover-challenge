// Loopover board simulator that evaluates moves in Programmer Notation.
// This version is a port from spdskatr's Python version, which can be found on GitHub:
// https://github.com/Loopover/LoopoverChallenge/blob/master/evaluator_oo.py
//
// Original author: spdskatr.
// Author: Martín "Netux" Rodriguez (https://me.netux.site).
// License: MIT.
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"unicode"
)

/* Programmer's Notation for a move
move 									  = amount axis index [reverse-index-indicator]

axis 									  = horizontal-axis | vertical-axis
index										= number
amount 								  = [ backwards-indicator ] number
backwards-indicator 		= -
reverse-index-indicator = '

horizontal-axis					= "R" | "r"
vertical-axis					  = "C" | "c"

number									 = { digit }
positive-number 		     = positive-digit [ { digit } ]
digit								   	 = "0" | positive-digit
positive-digit					 = "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
*/

// Board is a width*height Loopover board.
type Board [][]int

// NewBoard creates a new Board with the given dimensions.
func NewBoard(width, height int) (Board, error) {
	if width <= 1 {
		return nil, fmt.Errorf("board width must be greater than 1")
	}
	if height <= 1 {
		return nil, fmt.Errorf("board height must be greater than 1")
	}

	b := make(Board, width, width)
	for x := range b {
		b[x] = make([]int, height, height)

		for y := range b[x] {
			b.resetTile(x, y)
		}
	}

	return b, nil
}

// defaultTileValue returns the default value of a tile according to its coordinate.
func (b *Board) defaultTileValue(x, y int) int {
	return x + y*b.Width() + 1
}

// resetTile sets the tile's value to its default.
func (b *Board) resetTile(x, y int) {
	(*b)[x][y] = b.defaultTileValue(x, y)
}

// Width returns the width of the board.
func (b *Board) Width() int {
	return cap(*b)
}

// Height returns the height of the board by looking at the capacity of it's first column.
func (b *Board) Height() int {
	return cap((*b)[0])
}

// Reset resets the board by setting all tiles in the default order.
func (b *Board) Reset() {
	for x := 0; x < b.Width(); x++ {
		for y := 0; y < b.Height(); y++ {
			b.resetTile(x, y)
		}
	}
}

// FastShuffle shuffles the board by going through all the tiles and swaping them with a different, random tile.
func (b *Board) FastShuffle() {
	board := *b

	for x1 := 0; x1 < b.Width(); x1++ {
		for y1 := 0; y1 < b.Height(); y1++ {
			var x2, y2 int
			for x2, y2 = x1, y1; x1 == x2 || y1 == y2; {
				// keep generating a random x2 and y2 if x1 is equal to x2 and y1 is equal to y2.
				// this is done to avoid keeping a number on the same tile.
				x2 = rand.Intn(b.Width())
				y2 = rand.Intn(b.Height())
			}

			board[x1][y1], board[x2][y2] = board[x2][y2], board[x1][y1]
		}
	}
}

// Shuffle shuffles the board by applying `iterations` anmount of Moves generated with random parameters. If `iterations` is less or equal to 0, b.Width() + b.Height() is used instead.
// While this might be slower with more iterations, it is more truthful to what a human would do if they were to shuffle manually.
func (b *Board) Shuffle(iterations int) int {
	if iterations == 0 {
		iterations = b.Width() + b.Height()
	}

	var moves int
	for moves = 0; moves < iterations; moves++ {
		var a Axis
		var max int
		if rand.Intn(2) == 0 {
			a = HorizontalAxis
			max = b.Width()
		} else {
			a = VerticalAxis
			max = b.Height()
		}

		b.MakeMove(&Move{
			Axis:   a,
			Index:  rand.Intn(max),
			Amount: rand.Intn(max-1) + 1,
		})
	}

	return moves
}

// IsSolved returns true if all tiles are in order.
// In other words, if for every (x, y) the tile at (x, y) equals x + y * b.Width() + 1.
func (b *Board) IsSolved() bool {
	for x := 0; x < b.Width(); x++ {
		for y := 0; y < b.Height(); y++ {
			if (*b)[x][y] != b.defaultTileValue(x, y) {
				return false
			}
		}
	}

	return true
}

// MakeMove modifies the Board by applying a move. A move can shift the contents of a column or a row forward or backwards.
func (b *Board) MakeMove(m *Move) int {
	board := *b
	amnt := Abs(m.Amount)
	forward := m.Amount > 0

	for i := 0; i < amnt; i++ {
		if m.Axis == HorizontalAxis {
			if forward {
				p := board[b.Width()-1][m.Index]
				for x := 0; x < b.Width(); x++ {
					board[x][m.Index], p = p, board[x][m.Index]
				}
			} else {
				p := board[0][m.Index]
				for x := b.Width() - 1; x >= 0; x-- {
					board[x][m.Index], p = p, board[x][m.Index]
				}
			}
		} else {
			if forward {
				p := board[m.Index][b.Height()-1]
				for y := 0; y < b.Height(); y++ {
					board[m.Index][y], p = p, board[m.Index][y]
				}
			} else {
				p := board[m.Index][0]
				for y := b.Height() - 1; y >= 0; y-- {
					board[m.Index][y], p = p, board[m.Index][y]
				}
			}
		}
	}

	return amnt
}

// Axis can either be Horizontal or Vertical.
type Axis bool

const (
	// HorizontalAxis goes from left to right or vice-versa.
	HorizontalAxis = iota == 0
	// VerticalAxis goes from top to bottom or vice-versa.
	VerticalAxis
)

// Move represents a parsed Move from input.
type Move struct {
	Axis   Axis
	Index  int
	Amount int
}

// ParseMove creates a parsed Move from an input string in Programmer's Notation.
func ParseMove(input string, board *Board) (*Move, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	// parse axis by getting it's index.
	ai := -1
	for i, c := range input {
		if unicode.IsLetter(c) {
			ai = i
			break
		}
	}

	if ai == -1 {
		return nil, fmt.Errorf("no move character in move %q", input)
	}

	var axis Axis
	switch c := []rune(input)[ai]; unicode.ToLower(c) {
	default:
		return nil, fmt.Errorf("invalid move character %c in move %q", c, input)
	case 'r':
		axis = HorizontalAxis
	case 'c':
		axis = VerticalAxis
	}

	// parse amount.
	amountStr := input[:ai]
	if len(amountStr) == 0 {
		return nil, fmt.Errorf("missing amount in move %q", input)
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return nil, fmt.Errorf("invalid number for amount in move %q", input)
	}

	// check if amount is in bounds.
	if amount == 0 {
		return nil, fmt.Errorf("amount cannot be 0 in move %q", input)
	}

	// parse index.
	indexStr := input[ai+1:]
	if len(indexStr) == 0 {
		return nil, fmt.Errorf("missing amount in move %q", input)
	}

	var reverseIndex bool
	if input[len(input)-1] == '\'' {
		indexStr = indexStr[:len(indexStr)-1]
		reverseIndex = true
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return nil, fmt.Errorf("invalid number for index in move %q", input)
	}

	if reverseIndex {
		if axis == HorizontalAxis {
			index = board.Width() - 1 - index
		} else {
			index = board.Height() - 1 - index
		}
	}

	// check if index is in bounds.
	if index < 0 {
		return nil, fmt.Errorf("amount must be greater or equal to 0 in move %q", input)
	}

	var max int
	if axis == HorizontalAxis {
		max = board.Width()
	} else {
		max = board.Height()
	}

	if index >= max {
		return nil, fmt.Errorf("amount must be lesser to %d in move %q", max, input)
	}

	return &Move{
		Axis:   axis,
		Amount: amount,
		Index:  index,
	}, nil
}

// Abs returns the absolute value of the integer `a`.
func Abs(a int) int {
	if a > 0 {
		return a
	}

	return -a
}

// ParseTwoDimensions converts a string in the format of "WIDTHxHEIGHT", "widthXheight", "width*height" into two variables width and height.
func ParseTwoDimensions(input string) (width, height int, err error) {
	if len(input) == 0 {
		err = fmt.Errorf("empty input")
		return
	}

	i := strings.IndexAny(input, "xX*")

	if i == -1 {
		err = fmt.Errorf("%q doesn't seem to be a valid dimension", input)
		return
	}

	width, err = strconv.Atoi(input[:i])
	if err != nil {
		return
	}

	height, err = strconv.Atoi(input[i+1:])
	if err != nil {
		return
	}

	return
}

// SprintBoard formats the board into a grid of rows and columns.
func SprintBoard(b *Board) string {
	var r string
	pad := strconv.Itoa(len(strconv.Itoa(b.Width() * b.Height())))

	for y := 0; y < b.Height(); y++ {
		for x := 0; x < b.Width(); x++ {
			r += fmt.Sprintf(" %"+pad+"d", (*b)[x][y])
		}

		if y != b.Height()-1 {
			r += "\n"
		}
	}

	return r
}

// ScanShuffle scans user input to answer certain questions and execute either a fast or a normal shuffle.
func ScanShuffle(b *Board, scanner *bufio.Scanner) {
	var done bool

	fmt.Print("Fast shuffle? [Y/n]: ")
	for !done && scanner.Scan() {
		s := strings.ToLower(scanner.Text())
		if s != "n" {
			b.FastShuffle()
			fmt.Println("Fast shuffled board")

			done = true
			break
		}

		var iters int
		fmt.Print("How many iterations? 0 uses the default number of iterations: ")
		for scanner.Scan() {
			s = scanner.Text()

			// empty input means 0 iterations.
			if len(s) != 0 {
				iters, err := strconv.Atoi(s)
				if err != nil {
					fmt.Print("Invalid number, try again: ")
					continue
				}

				if iters < 0 {
					iters = 0
				}
			}

			break
		}

		finalIters := b.Shuffle(iters)
		fmt.Printf("Shuffled board with %d iterations\n", finalIters)

		done = true
	}
}

func main() {
	var b Board
	scanner := bufio.NewScanner(os.Stdin)

	// scan board size.
	fmt.Print("Input board size (default is 5x5): ")
	for scanner.Scan() {
		var w, h int
		var err error

		s := scanner.Text()

		if s == "" {
			w = 5
			h = 5
		} else {
			var err error
			w, h, err = ParseTwoDimensions(s)
			if err != nil {
				fmt.Printf("Invalid size (%s), try again: ", err)
				continue
			}
		}

		b, err = NewBoard(w, h)
		if err != nil {
			fmt.Printf("Error creating board (%s), try again: ", err)
			continue
		}

		break
	}

	n := 0

	// game loop.
	for {
		fmt.Println()

		// present board state.
		fmt.Println("Board state:")
		fmt.Println(SprintBoard(&b))

		fmt.Printf("%d moves so far\n", n)

		if b.IsSolved() {
			fmt.Println("Solved")
		}

		// scan for moves.
		fmt.Print("Move: ")
		for scanner.Scan() {
			switch s := strings.ToLower(scanner.Text()); s {
			case "shuffle":
				ScanShuffle(&b, scanner)
			case "reset":
				b.Reset()
				n = 0
				fmt.Println("Board reset")
			default:
				m, err := ParseMove(s, &b)
				if err != nil {
					fmt.Printf("Invalid move (%s), try again: ", err)
					continue
				}

				n += b.MakeMove(m)
			}

			break
		}
	}
}
