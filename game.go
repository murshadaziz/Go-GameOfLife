package main

import (
	"fmt"
	"math/rand/v2"
)

const (
	rows    int  = 25
	columns int  = 80
	dead    rune = '.'
	alive   rune = 'o'
)

type Universe [][]bool

type Game struct {
	uni1 Universe
	uni2 Universe
	

}

func NewUniverse() Universe {
	buf := make([]bool, rows*columns)
	var uni Universe = make(Universe, rows)
	for i := range uni {
		uni[i] = buf[i*columns : (i+1)*columns]
	}
	return uni
}

func (uni Universe) show() {
	for i := range uni {
		for j := range uni[i] {
			if uni[i][j] == false {
				fmt.Printf("%c", dead)
			} else {
				fmt.Printf("%c", alive)
			}
		}
		fmt.Println()
	}

}

func (uni Universe) seed() {
	for range rows * columns / 4 {
		uni[rand.IntN(rows)][rand.IntN(columns)] = true
	}
}

func (uni Universe) alive(i, j int) bool {
	return uni[(i+rows)%rows][(j+columns)%columns]
}

func (uni Universe) count(i, j int) int {
	count := 0
	for x := i - 1; x <= i+1; x++ {
		for y := j - 1; y <= j+1; y++ {
			if x == i && y == j {
				continue
			}
			if uni.alive(x, y) {
				count++
			}
		}
	}
	return count
}

func step(uni1, uni2 *Universe) {
	for i := range *uni1 {
		for j := range (*uni1)[i] {
			count := (*uni1).count(i, j)
			if (*uni1)[i][j] == true && (count < 2 || count > 3) {
				(*uni2)[i][j] = false
			} else if (*uni1)[i][j] == false && count == 3 {
				(*uni2)[i][j] = true
			} else {
				(*uni2)[i][j] = (*uni1)[i][j]
			}
		}
	}
	*uni1, *uni2 = *uni2, *uni1
}

func main() {
	var uni1 Universe = NewUniverse()
	var uni2 Universe = NewUniverse()
	uni1[0][0] = true
	uni1[0][1] = true
	uni1[0][2] = true
	uni1[1][0] = true
	uni1[1][1] = true
	uni1[1][2] = true
	uni1.show()
	fmt.Println("\x0c")
	step(&uni1, &uni2)
	uni1.show()
}
