// input.go
package main

import (
	"github.com/eiannone/keyboard"
)

func listenForInput(keys chan<- keyEvent, quit chan<- struct{}) {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return
		}
		if char == 'q' || key == keyboard.KeyCtrlC {
			close(quit)
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
				break welcomeloop
			case keychar.char == 'r':
				G.seed()
				G.render(100, quit)
				break welcomeloop
			default:
				continue
			}
		}
	}
}
