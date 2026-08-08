// input.go
package main

import "os"

func listenForInput(keys chan<- rune) {
	for {
		var b []byte = make([]byte, 1)
		os.Stdin.Read(b)
		keys <- rune(b[0])
	}
}
func gameLoop(keys <-chan rune) {
	for {
		key := <-keys

		switch key {
		case 'q':
			return

		case ' ':
			// advance simulation

		case 's':
			// seed / something

		case 'p':
			// pause
		}
	}
}
