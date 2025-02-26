package main

import(
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// os.Args[0] is the script name, arguments start from os.Args[1]
	if len(os.Args) < 2 {
		fmt.Println("No arguments provided!")
		return
	}
	cmd := strings.Split(os.Args[1], " ")
	filepath.Walk("/termfiles/.packages", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == os.Args[1] + ".go" {
			output, err := exec.Command("go", "run", path, cmd[1]).CombinedOutput()
			if err != nil {
				fmt.Println("Error:", err)
			}
			  
			// Print script output
			fmt.Println(string(output))
			fmt.Println("Found:", path)
		}
		return nil
	})
	// Print all arguments
	fmt.Println("Arguments received:")
	for i, arg := range os.Args[1:] {
		fmt.Printf("Arg %d: %s\n", i+1, arg)
	}
}