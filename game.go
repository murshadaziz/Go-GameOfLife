package main

import (
	"fmt"
	"math/rand"
)

const (
	rows    int  = 25
	columns int  = 40
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
	for range 25 {
		uni[rand.Intn(rows)][rand.Intn(columns)] = true
	}
}

func main() {
	var uni Universe = NewUniverse()
	uni.seed()
	uni.show()
}
