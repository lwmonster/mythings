package main

import (
    "bufio"
    "os"
    "fmt"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter text: ")
    text, _ := reader.ReadString('\n')
    fmt.Println(text)

    fmt.Println("Enter text: ")
    text2 := ""
    fmt.Scanln(&text2)
    fmt.Println(text2)

}
