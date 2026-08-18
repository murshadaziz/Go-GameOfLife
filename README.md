# Go Game of Life

A terminal implementation of Conway's Game of Life, written in Go. Renders directly to the terminal using ANSI escape sequences, supports live keyboard input, and includes an interactive editor for drawing custom starting patterns.

## Features

- Full-screen terminal rendering using the alternate screen buffer
- Automatically sizes the grid to fit the current terminal dimensions
- Toroidal (wrap-around) universe: cells on the edges neighbor cells on the opposite edge
- Interactive pattern editor with a movable cursor
- Random seeding
- Pause, resume, and clear the grid at any time
- Clean shutdown on Ctrl+C or SIGTERM, restoring the terminal state

## Controls

**Welcome screen**

| Key | Action |
|-----|--------|
| `d` | Draw a custom starting pattern |
| `r` | Seed the grid randomly |
| `Ctrl+C` | Quit |

**Editor mode**

| Key | Action |
|-----|--------|
| `w` `a` `s` `d` | Move the cursor |
| `Space` | Toggle the cell under the cursor |
| `r` | Reseed randomly |
| `Enter` | Confirm and start the simulation |

**Simulation**

| Key | Action |
|-----|--------|
| `p` | Pause / resume |
| `r` | Reseed randomly (while paused) |
| `c` | Clear the grid (while paused) |
| `d` | Open the editor (while paused) |
| `Ctrl+C` | Quit |

## Requirements

- Go 1.26.5 or later

## Installation

```bash
git clone https://github.com/murshadaziz/Go-GameOfLife.git
cd Go-GameOfLife
go mod tidy
```

## Usage

```bash
go run .
```

Or build a binary:

```bash
go build -o game-of-life
./game-of-life
```

## Project Structure

```
.
├── game.go   # Game state, rendering, editor, and simulation loop
├── input.go  # Keyboard input listener and welcome screen loop
├── go.mod
└── go.sum
```

## Dependencies

- [github.com/eiannone/keyboard](https://github.com/eiannone/keyboard) - keyboard input handling
- [golang.org/x/term](https://pkg.go.dev/golang.org/x/term) - terminal size detection

## Demo

*(video coming soon)*

## License

No license specified.