package main


import (
    "fmt"
    "sync"
)


func main() {
    var once sync.Once
    onceBody := func() {
        fmt.Println("Only once")
    }

    done := make(chan int)

    for i := 0; i < 10; i++ {
        go func() {
            once.Do(onceBody)
            done <- 4
        }()
    }
    for i := 0; i < 9; i++ {
        fmt.Println( <-done );
    }
}

