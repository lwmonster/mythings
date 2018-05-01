package main

import (
    "fmt"
    "math/rand"
    "time"
)

func printRandom () {
    for i := 0; i < 10; i++  {
        fmt.Println(rand.Intn(100))
    }
}


func main() {
    rand.Seed(time.Now().UnixNano())
    printRandom()
}

