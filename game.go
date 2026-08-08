package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

const (
	rows         int  = 25
	columns      int  = 80
	dead         rune = '.'
	alive        rune = 'o'
	altScreenOn       = "\x1b[?1049h"
	altScreenOff      = "\x1b[?1049l"
	cursorHide        = "\x1b[?25l"
	cursorShow        = "\x1b[?25h"
	clearScreen       = "\x1b[2J"
	cursorHome        = "\x1b[H"
)

type Universe [][]bool

type Game struct {
	uni1     Universe
	uni2     Universe
	Rows     int
	Columns  int
	isPaused bool
	buffer   *bufio.Writer
}

func (G *Game) init() {
	buf1 := make([]bool, G.Rows*G.Columns)
	G.uni1 = make(Universe, G.Rows)
	G.buffer = bufio.NewWriter(os.Stdout)

	for i := range G.uni1 {
		G.uni1[i] = buf1[i*G.Columns : (i+1)*G.Columns]
	}

	buf2 := make([]bool, G.Rows*G.Columns)
	G.uni2 = make(Universe, G.Rows)

	for i := range G.uni2 {
		G.uni2[i] = buf2[i*G.Columns : (i+1)*G.Columns]
	}
}

func (G *Game) seed() {
	for range G.Rows * G.Columns / 4 {
		G.uni1[rand.IntN(G.Rows)][rand.IntN(G.Columns)] = true
	}
}

func (G *Game) alive(i, j int) bool {
	return G.uni1[(i+G.Rows)%G.Rows][(j+G.Columns)%G.Columns]
}

func (G *Game) count(i, j int) int {
	count := 0
	for x := i - 1; x <= i+1; x++ {
		for y := j - 1; y <= j+1; y++ {
			if x == i && y == j {
				continue
			}
			if G.alive(x, y) {
				count++
			}
		}
	}
	return count
}

func (G *Game) step() {
	for i := range G.uni1 {
		for j := range G.uni1[i] {
			count := G.count(i, j)
			if G.uni1[i][j] == true && (count < 2 || count > 3) {
				G.uni2[i][j] = false
			} else if G.uni1[i][j] == false && count == 3 {
				G.uni2[i][j] = true
			} else {
				G.uni2[i][j] = G.uni1[i][j]
			}
		}
	}
	G.uni1, G.uni2 = G.uni2, G.uni1
}

func (G *Game) draw() {
	G.buffer.Write([]byte(cursorHome))
	for i := range G.uni1 {
		for j := range G.uni1[i] {
			if G.uni1[i][j] {
				G.buffer.WriteByte(byte(alive))
			} else {
				G.buffer.WriteByte(byte(dead))
			}
		}
		G.buffer.WriteByte('\n')
	}
	G.buffer.Flush()
}

func (G *Game) render() {
	fmt.Print(altScreenOn + cursorHide + clearScreen)
	defer fmt.Print(cursorShow + altScreenOff)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		G.step()
		G.draw()
	}
}

func main() {
	game := &Game{
		Rows:    rows,
		Columns: columns,
	}
	game.init()
	game.seed()
	game.render()
}
