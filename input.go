// input.go
package main

import (
	"github.com/eiannone/keyboard"
)

func listenForInput(keys chan<- keyEvent, closeQuit func()) {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return
		}
		if key == keyboard.KeyCtrlC {
			closeQuit()
			return
		}
		keys <- keyEvent{char: char, key: key}
	}
}
func (G *Game) WelcomeLoop(keys <-chan keyEvent, quit <-chan struct{}) {
welcomeloop:
	for {
		select {
		case <-quit:
			return
		case keychar := <-keys:
			switch {
			case keychar.char == 'd':
				G.runEditor(keys, quit)
				G.render(100, quit, keys)
				break welcomeloop
			case keychar.char == 'r':
				G.seed()
				G.render(100, quit, keys)
				break welcomeloop
			default:
				continue
			}
		}
	}
}
