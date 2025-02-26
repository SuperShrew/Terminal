package main

import(
	"fmt"
	"os"
)

func main() {
	// os.Args[0] is the script name, arguments start from os.Args[1]
	if len(os.Args) < 2 {
		fmt.Println("No arguments provided!")
		return
	}

	// Print all arguments
	fmt.Println("Arguments received:")
	for i, arg := range os.Args[1:] {
		fmt.Printf("Arg %d: %s\n", i+1, arg)
	}
}