package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/term"

	"github.com/eiannone/keyboard"
)

const (
	dead  byte = ' '
	alive byte = 'o'
	// ANSI escape codes for terminal control sequences
	altScreenOn  = "\x1b[?1049h"
	altScreenOff = "\x1b[?1049l"
	cursorHide   = "\x1b[?25l"
	cursorShow   = "\x1b[?25h"
	clearScreen  = "\x1b[2J"
	cursorHome   = "\x1b[H"
	cursorMoveTo = "\x1b[%d;%dH"
	clearLine    = "\x1b[K"
)

// keyEvent struct represents a keyboard event with the character and key pressed.
type keyEvent struct {
	char rune
	key  keyboard.Key
}

type Universe [][]bool // A 2D slice of booleans representing the game grid

// Game struct encapsulates the state and behavior of the game
type Game struct {
	uni1      Universe
	uni2      Universe
	Rows      int
	Columns   int
	isPaused  bool
	buffer    *bufio.Writer
	cursorRow int
	cursorCol int
}

// setTerminalSize retrieves the current terminal size and updates rows and columns in the Game struct
func (G *Game) setTerminalSize() {
	columns, rows, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Could not get terminal size:", err)
		return
	}

	G.Columns = columns - 2
	G.Rows = rows - 2
}

// init initializes the game by creating two universes and a buffered writer for output.
func (G *Game) init() {
	G.setTerminalSize()
	buf1 := make([]bool, G.Rows*G.Columns)
	G.uni1 = make(Universe, G.Rows)
	/* Create a buffered writer that accumulates output before writing it to stdout.
	the NewWriter function creates a new buffered writer that writes to the specified io.Writer (in this case, os.Stdout).
	NewWriter accepts anthing that implements the io.Writer interface(has method Write([]byte)), and returns a pointer to a new bufio.Writer.
	*/
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

// Randomly seeds the first universe with live cells.
func (G *Game) seed() {
	for range G.Rows * G.Columns / 4 {
		G.uni1[rand.IntN(G.Rows)][rand.IntN(G.Columns)] = true
	}
}

// alive checks if the cell at position (i, j) is alive, wrapping around the edges of the universe.
func (G *Game) alive(i, j int) bool {
	return G.uni1[(i+G.Rows)%G.Rows][(j+G.Columns)%G.Columns]
}

// count counts the number of alive neighbors for the cell at position (i, j).
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
	G.uni1, G.uni2 = G.uni2, G.uni1 //Swap the two universes
}

// draw renders the current state of the universe to the terminal.
func (G *Game) draw() {
	//This writer accumulates data in an internal buffer and writes it to the underlying writer (os.Stdout) when the buffer is full or when Flush() is called.
	G.buffer.Write([]byte(cursorHome))
	for i := range G.uni1 {
		for j := range G.uni1[i] {
			if G.uni1[i][j] {
				G.buffer.WriteByte(alive)
			} else {
				G.buffer.WriteByte(dead)
			}
		}
		G.buffer.WriteByte('\n')
	}
	G.drawFooter("p pause   c clear   d draw   r random   ctrl+c quit")
	G.buffer.Flush() // Flushes the buffer, writing any accumulated data to the underlying writer (os.Stdout).
}

func (G *Game) drawFooter(text string) {
	pad := spaces((G.Columns - len([]rune(text))) / 2)
	G.buffer.WriteString(pad)
	G.buffer.WriteString(text)
	G.buffer.WriteString(clearLine)
	G.buffer.WriteByte('\n')
}

// render sets up the game loop, updating and drawing the universe at a fixed frame rate.
func (G *Game) render(miliseconds int, quit <-chan struct{}, keys <-chan keyEvent) {
	// ANSI escape codes to switch to the alternate screen buffer, hide the cursor, and clear the screen.
	fmt.Print(altScreenOn + cursorHide + clearScreen)
	// The defer statement ensures that the cursor is shown and the alternate screen buffer is turned off when the function returns, even if an error occurs or the function exits early.
	defer fmt.Print(cursorShow + altScreenOff)
	// Defines the frame rate for the game loop.
	ticker := time.NewTicker(time.Duration(miliseconds) * time.Millisecond)

	// Stops the the ticker when the function returns, releasing any resources associated with it.
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !G.isPaused {
				G.step()
				G.draw()
			}
		case <-quit:
			return

		case ev := <-keys:
			switch {
			case ev.char == 'p':
				G.isPaused = !G.isPaused
				G.draw()
				if G.isPaused {
					G.drawPaused()
				}
			case ev.char == 'r' && G.isPaused:
				G.seed()
				G.draw()
			case ev.char == 'c' && G.isPaused:
				G.clearGrid()
				G.draw()
			case ev.char == 'd' && G.isPaused:
				G.runEditor(keys, quit)
				G.draw()
			}
		}
	}
}

// Switches the state of the cell at position (i, j) between alive and dead.
func (G *Game) toggleCell(i, j int) {
	G.uni1[i][j] = !G.uni1[i][j]
}

// clearGrid sets all cells in both universes to dead.
func (g *Game) clearGrid() {
	for i := range g.uni1 {
		for j := range g.uni1[i] {
			g.uni1[i][j] = false
			g.uni2[i][j] = false
		}
	}
}
func spaces(n int) string {
	if n < 0 {
		n = 0
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}

// drawWelcome draws the title screen and the menu: d = draw, r = random, q = quit.
func (G *Game) drawWelcome() {
	G.buffer.WriteString(altScreenOn + clearScreen + cursorHome)

	title := "GAME OF LIFE"
	pad := spaces((G.Columns - len(title)) / 2)

	G.buffer.WriteString("\n\n")
	G.buffer.WriteString(pad)
	G.buffer.WriteString(title)
	G.buffer.WriteString("\n\n")
	G.buffer.WriteString(pad)
	G.buffer.WriteString("d - draw your own pattern\n")
	G.buffer.WriteString(pad)
	G.buffer.WriteString("r - random seed\n")
	G.buffer.WriteString(pad)
	G.buffer.WriteString("ctrl+c - quit\n")
	G.buffer.Flush()
}

// drawWithCursor renders the grid like draw(), but marks the cursor's
// current position so you can see where space will toggle a cell.
func (G *Game) drawWithCursor() {
	G.buffer.Write([]byte(cursorHome))
	for i := range G.uni1 {
		for j := range G.uni1[i] {
			switch {
			case i == G.cursorRow && j == G.cursorCol:
				G.buffer.WriteByte('I') // Mark the cursor position with 'I'
			case G.uni1[i][j]:
				G.buffer.WriteByte(alive)
			default:
				G.buffer.WriteByte(dead)
			}
		}
		G.buffer.WriteByte('\n')
	}
	G.drawFooter("wasd move   space toggle   enter confirm")
	G.buffer.Flush()
}

// drawPaused displays a message in the center of the screen when the game is paused.
func (G *Game) drawPaused() {
	msg := " PAUSED: p resume   c clear   d draw   r random   ctrl+c quit "
	col := (G.Columns - len(msg)) / 2
	row := G.Rows / 2
	fmt.Fprintf(G.buffer, cursorMoveTo+"%s", row, col+1, msg)
	G.buffer.Flush()
}

// runEditor lets the user move a cursor around the grid with h/j/k/l,
// toggle the cell under it with space, reseed with 'r', and hand off to
// the simulation with Enter. 'q'/Ctrl+C are caught upstream in
// listenForInput, which closes quit — so this func doesn't need to check
// for them itself.
func (G *Game) runEditor(keys <-chan keyEvent, quit <-chan struct{}) {
	G.cursorRow, G.cursorCol = G.Rows/2, G.Columns/2
	G.drawWithCursor()

	for {
		select {
		case <-quit:
			return
		case key := <-keys:
			switch {
			case key.char == 'a':
				if G.cursorCol > 0 {
					G.cursorCol--
				}
			case key.char == 'd':
				if G.cursorCol < G.Columns-1 {
					G.cursorCol++
				}
			case key.char == 'w':
				if G.cursorRow > 0 {
					G.cursorRow--
				}
			case key.char == 's':
				if G.cursorRow < G.Rows-1 {
					G.cursorRow++
				}
			case key.key == keyboard.KeySpace:
				G.toggleCell(G.cursorRow, G.cursorCol)
			case key.char == 'r':
				G.seed()
			case key.key == keyboard.KeyEnter:
				// enter: done drawing
				return
			default:
				continue
			}
			G.drawWithCursor()
		}
	}
}

func main() {
	game := &Game{}
	keyboard.Open()
	defer keyboard.Close()
	game.init()
	keys := make(chan keyEvent)
	quit := make(chan struct{})
	var quitOnce sync.Once // Ensures that the quit channel is closed only once

	// closeQuit is a helper function to close the quit channel safely.
	closeQuit := func() {
		quitOnce.Do(func() {
			close(quit)
		})
	}

	// Catch Ctrl+C (SIGINT) and kill/SIGTERM.
	sigCh := make(chan os.Signal, 1)                    // Create a channel to receive OS signals with a buffer size of 1
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM) // Catches the interrupt signal (Ctrl+C) and the termination signal (SIGTERM) and sends them to sigCh
	defer signal.Stop(sigCh)                            // Stops the signal notification when the main function exits, preventing any further signals from being sent to sigCh

	// Listen for OS signals.
	go func() {
		<-sigCh
		fmt.Println(clearScreen + cursorShow + altScreenOff) // Clear the screen and show the cursor before exiting
		closeQuit()                                          // Close the quit channel when a signal is received, signaling the game to exit
	}()
	game.drawWelcome()
	go listenForInput(keys, closeQuit) // Start listening for keyboard input in a separate goroutine
	game.WelcomeLoop(keys, quit)

}
