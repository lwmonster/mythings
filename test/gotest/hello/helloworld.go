package main

import "fmt"


type Op struct {
    Cmd string `json:"cmd"`
}

func main() {
    fmt.Println("Hello World！！")

    op := Op{}

    fmt.Println(op)

}
