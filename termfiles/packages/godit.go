package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
    "strings"
	"io/ioutil"
)

func Edit(preFilled string) string {
	fd := int(os.Stdin.Fd())

	// Save terminal state
	var oldState syscall.Termios
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TCGETS, uintptr(unsafe.Pointer(&oldState)))

	// Set new terminal state
	newState := oldState
	newState.Lflag &^= syscall.ICANON | syscall.ECHO
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TCSETS, uintptr(unsafe.Pointer(&newState)))
	defer syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TCSETS, uintptr(unsafe.Pointer(&oldState)))

	// Print the initial prompt
	fmt.Println("godit v1.0")

	// Convert the pre-filled string into lines
	lines := strings.Split(preFilled, "\n")
	if len(lines) == 0 {
		lines = append(lines, "")
	}
	currentLine := len(lines) - 1
	cursorX := len(lines[currentLine])

	for {
		var buf [1]byte
		os.Stdin.Read(buf[:])
		key := buf[0]

		// Ctrl+X to save and exit
		if key == 24 { // ASCII code for Ctrl+X
			break
		}

		// Handle backspace
		if key == 8 || key == 127 {
			if cursorX > 0 {
				lines[currentLine] = lines[currentLine][:cursorX-1] + lines[currentLine][cursorX:]
				cursorX--
			} else if currentLine > 0 { // Merge with previous line
				cursorX = len(lines[currentLine-1])
				lines[currentLine-1] += lines[currentLine]
				lines = append(lines[:currentLine], lines[currentLine+1:]...)
				currentLine--
			}
		} else if key == 27 {
			// Handle arrow keys
			var seq [2]byte
			os.Stdin.Read(seq[:])
			if seq[0] == '[' {
				if seq[1] == 'A' && currentLine > 0 { // Move up one line
					currentLine--
					cursorX = min(cursorX, len(lines[currentLine]))
				}
				if seq[1] == 'B' && currentLine < len(lines)-1 { // Move down one line
					currentLine++
					cursorX = min(cursorX, len(lines[currentLine]))
				}
				if seq[1] == 'D' && cursorX > 0 { // Move left
					cursorX--
				}
				if seq[1] == 'C' && cursorX < len(lines[currentLine]) { // Move right
					cursorX++
				}
			}
		} else if key == '\n' {
			// Handle enter key (new line)
			newLine := lines[currentLine][cursorX:]
			lines[currentLine] = lines[currentLine][:cursorX]
			lines = append(lines[:currentLine+1], append([]string{newLine}, lines[currentLine+1:]...)...)
			currentLine++
			cursorX = 0
		} else {
			// Insert character at cursor position
			lines[currentLine] = lines[currentLine][:cursorX] + string(key) + lines[currentLine][cursorX:]
			cursorX++
		}

		// Clear screen and reprint with correct cursor positioning
		fmt.Print("\033[H\033[2J") // Clear screen
		fmt.Println("godit v1.0")
		for i, line := range lines {
			if i == currentLine {
				fmt.Print("→ ") // Indicate the current line
			}
			fmt.Println(line)
		}
		fmt.Printf("\033[%d;%dH", currentLine+2, cursorX+3) // Move cursor to correct position
	}

	fmt.Println() // New line after exiting
	return strings.Join(lines, "\n")
}

// Helper function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
	Edit(string(data))
}