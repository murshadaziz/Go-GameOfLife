package main

import (
	"fmt"
	"math/rand"
)

const (
	rows    int  = 25
	columns int  = 80
	dead    rune = '.'
	alive   rune = 'o'
)

type Universe [][]bool

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
		uni[rand.Intn(rows)][rand.Intn(columns)] = true
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

func main() {
	var uni Universe = NewUniverse()
	uni.seed()
	uni.show()
	var c int = uni.count(1, 2)
	fmt.Println(c)
}
