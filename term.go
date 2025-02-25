package main

import (
  "fmt"
  "strings"
  "os"
  "os/user"
  "bufio"
)

func prompt(user string) string {
  dir, err := os.Getwd()
  if err != nil {
    fmt.Println("Error:", err)
    return ""
  }
  return dir + " - " + user + "@go-term $ "
}

func excecute(command string) int {
  command = strings.TrimSpace(command)
  dir, err := os.Getwd()
  if err != nil {
    fmt.Println("Error:", err, dir)
    return 0
  }
  //fmt.Println(dir)
  //fmt.Println(command)
  cmd := strings.Split(command, " ")
  //fmt.Printf("Raw input: %q\n", cmd[0])
  if cmd[0] == "cd" {
    err = os.Chdir(cmd[1])
    if err != nil {
      fmt.Println(err)
      return 0
    }
    return 1
  } else if cmd[0] == "ls" {
    files, err := os.ReadDir(dir)
    if err != nil {
      fmt.Println("Error reading directory:", err)
      return 0
    }
    for _, file := range files {
      fmt.Println(file.Name())
    }
  }
  return 0
}

func main() {
  var code int
  var lastCommand string
  var command string
  currentUser, err := user.Current()
  if err != nil {
    fmt.Println(err)
  }
  scanner := bufio.NewScanner(os.Stdin)
  for {
    fmt.Print(prompt(currentUser.Username))
    scanner.Scan()
    command = scanner.Text()
    code = excecute(command)
    lastCommand = command
    //fmt.Println("exit code:", code)
  }
}

