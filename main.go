package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const COLS int = 60
const ROWS int = 30

func main() {
	// HIDE CURSOR
	fmt.Print("\033[?25l")
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run() // the same as -icanon
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()              // do not display entered characters on the screen

	var key string
	var char = make(chan string)

	go func(channel chan string) {
		var c []byte = make([]byte, 1)
		for {
			os.Stdin.Read(c)
			channel <- string(c)
		}
	}(char)

	var x, y = make([]int, 1000), make([]int, 1000)
	var quit bool = false

	for !quit {
		RenderTable()
		// MOVE CURSOR BACK TO TOP
		fmt.Printf("\033[%dA", ROWS+2)
		var head, tail int = 0, 0
		x[head] = COLS / 2
		y[head] = ROWS / 2
		var gameOver bool = false
		var xdir, ydir int = 1, 0
		var applex, appley int = -1, 0

		for !quit && !gameOver {
			if applex < 0 {
				// CREATE NEW APPLE
				applex = rand.Intn(COLS)
				appley = rand.Intn(ROWS)

				for i := tail; i != head; i = (i + 1) % 1000 {
					if x[i] == applex && y[i] == appley {
						applex = -1
					}
				}

				if applex >= 0 {
					// DRAW APPLE
					fmt.Printf("\033[%dB\033[%dC❤", appley+1, applex+1)
					fmt.Printf("\033[%dF", appley+1)
				}
			}

			// CLEAR SNAKE TAIL
			fmt.Printf("\033[%dB\033[%dC·", y[tail]+1, x[tail]+1)
			fmt.Printf("\033[%dF", y[tail]+1)
			if x[head] == applex && y[head] == appley {
				applex = -1
				fmt.Print("\a") // THE BELL
			} else {
				tail = (tail + 1) % 1000
			}

			var newhead int = (head + 1) % 1000
			x[newhead] = (x[head] + xdir + COLS) % COLS
			y[newhead] = (y[head] + ydir + ROWS) % ROWS
			head = newhead

			for i := tail; i != head; i = (i + 1) % 1000 {
				if x[i] == x[head] && y[i] == y[head] {
					gameOver = true
				}
			}

			/// DRAW HEAD
			fmt.Printf("\033[%dB\033[%dC▓", y[head]+1, x[head]+1)
			fmt.Printf("\033[%dF", y[head]+1)

			time.Sleep(time.Microsecond * 84000)

			/// READ KEYBOARD
			select {
			case key = <-char:
				if key == "q" {
					quit = true
					exec.Command("stty", "-F", "/dev/tty", "echo").Run()
				} else if key == "a" && xdir != 1 {
					xdir = -1
					ydir = 0
				} else if key == "d" && xdir != -1 {
					xdir = 1
					ydir = 0
				} else if key == "s" && ydir != -1 {
					xdir = 0
					ydir = 1
				} else if key == "w" && ydir != 1 {
					xdir = 0
					ydir = -1
				}
			default:
			}
			/// END READ KEYBOARD
		}

		if !quit {
			/// SHOW GAMEOVER
			fmt.Printf("\033[%dB\033[%dC Game Over! ", ROWS/2, COLS/2-5)
			fmt.Printf("\033[%dF", ROWS/2)
		}
	}
	// SHOW CURSOR
	fmt.Print("\033[?25h")
}

func RenderTable() {
	fmt.Print("┌")
	for i := 0; i < COLS; i++ {
		fmt.Print("─")
	}
	fmt.Print("┐\n")
	////////////////////////////////////////
	for j := 0; j < ROWS; j++ {
		fmt.Print("│")
		for i := 0; i < COLS; i++ {
			fmt.Print("·")
		}
		fmt.Print("│\n")
	}
	////////////////////////////////////////
	fmt.Print("└")
	for i := 0; i < COLS; i++ {
		fmt.Print("─")
	}
	fmt.Print("┘\n")
}
