package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func readEditableInput(preFilled, lastString string) string {
    // Get the file descriptor for standard input (fd 0)
    fd := int(os.Stdin.Fd())

    // Save the current terminal settings
    var oldState syscall.Termios
    syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TCGETS, uintptr(unsafe.Pointer(&oldState)))

    // Create a copy of the terminal settings to modify
    newState := oldState

    // Disable canonical mode (for raw input processing) and echoing of input
    newState.Lflag &^= syscall.ICANON | syscall.ECHO

    // Set the new terminal settings
    syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TCSETS, uintptr(unsafe.Pointer(&newState)))

    // Ensure the terminal state is reset to the original state when the function exits
    defer syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TCSETS, uintptr(unsafe.Pointer(&oldState)))

    // Print the prompt with the pre-filled text
    fmt.Print("\r> " + preFilled)
    
    // Convert the pre-filled string into a slice of runes (to handle multi-byte characters)
    input := []rune(preFilled)

    // Set the initial cursor position to the end of the pre-filled input
    cursorPos := len(input)

    for {
        var buf [1]byte
        // Read one byte from stdin
        os.Stdin.Read(buf[:])
        key := buf[0]

        // If Enter (newline) is pressed, break the loop and return the input
        if key == '\n' {
            break
        }

        // Handle backspace (8 or 127 are ASCII codes for backspace)
        if key == 8 || key == 127 {
            if cursorPos > 0 {
                // Remove the character before the cursor
                input = append(input[:cursorPos-1], input[cursorPos:]...)
                cursorPos--
                // Print the updated input and move the cursor back one step
                fmt.Print("\r> " + string(input) + " \b")
                fmt.Print("\r> " + string(input[:cursorPos]))
            }
            continue
        }

        // Handle special key sequences, e.g., arrow keys
        if key == 27 {
            var seq [2]byte
            // Read the escape sequence (two more bytes)
            os.Stdin.Read(seq[:])
            if seq[0] == '[' {
                // Left arrow (move cursor left)
                if seq[1] == 'D' && cursorPos > 0 {
                    cursorPos--
                    fmt.Print("\b")
                }
                // Right arrow (move cursor right)
                if seq[1] == 'C' && cursorPos < len(input) {
                    fmt.Print(string(input[cursorPos]))
                    cursorPos++
                }
                // Up arrow (replace input with the last string)
                if seq[1] == 'A' {
                    input = []rune(lastString)
                    cursorPos = len(input)
                    // Display the last input string
                    fmt.Print("\r> " + string(input) + "  ")
                    fmt.Print("\r> " + string(input))
                }
            }
            continue
        }

        // Insert the character at the current cursor position
        if cursorPos == len(input) {
            input = append(input, rune(key))
        } else {
            input = append(input[:cursorPos], append([]rune{rune(key)}, input[cursorPos:]...)...)
        }

        cursorPos++
        // Print the updated input and reset the cursor position to the current position
        fmt.Print("\r> " + string(input) + " ")
        fmt.Print("\r> " + string(input[:cursorPos]))
    }

    // After exiting the loop, print a newline for clean output
    fmt.Println()
    
    // Return the final string that was entered
    return string(input)
}


func main() {
	lastCommand := "echo Hello, World!"
	input := readEditableInput("bob", lastCommand)
	fmt.Println("Final input:", input)
}
