package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func readEditableInput(preFilled, lastString string) string {
	fd := int(os.Stdin.Fd())
	var oldState syscall.Termios
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TCGETS, uintptr(unsafe.Pointer(&oldState)))
	newState := oldState
	newState.Lflag &^= syscall.ICANON | syscall.ECHO
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TCSETS, uintptr(unsafe.Pointer(&newState)))
	defer syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TCSETS, uintptr(unsafe.Pointer(&oldState)))

	fmt.Print("\r> " + preFilled)
	input := []rune(preFilled)
	cursorPos := len(input)

	for {
		var buf [1]byte
		os.Stdin.Read(buf[:])
		key := buf[0]

		if key == '\n' {
			break
		}

		if key == 8 || key == 127 {
			if cursorPos > 0 {
				input = append(input[:cursorPos-1], input[cursorPos:]...)
				cursorPos--
				fmt.Print("\r> " + string(input) + " \b")
				fmt.Print("\r> " + string(input[:cursorPos]))
			}
			continue
		}

		if key == 27 {
			var seq [2]byte
			os.Stdin.Read(seq[:])
			if seq[0] == '[' {
				if seq[1] == 'D' && cursorPos > 0 {
					cursorPos--
					fmt.Print("\b")
				}
				if seq[1] == 'C' && cursorPos < len(input) {
					fmt.Print(string(input[cursorPos]))
					cursorPos++
				}
				if seq[1] == 'A' { // Up arrow pressed
					input = []rune(lastString)
					cursorPos = len(input)
					fmt.Print("\r> " + string(input) + "  ")
					fmt.Print("\r> " + string(input))
				}
			}
			continue
		}

		if cursorPos == len(input) {
			input = append(input, rune(key))
		} else {
			input = append(input[:cursorPos], append([]rune{rune(key)}, input[cursorPos:]...)...)
		}

		cursorPos++
		fmt.Print("\r> " + string(input) + " ")
		fmt.Print("\r> " + string(input[:cursorPos]))
	}

	fmt.Println()
	return string(input)
}

func main() {
	lastCommand := "echo Hello, World!"
	input := readEditableInput("bob", lastCommand)
	fmt.Println("Final input:", input)
}
