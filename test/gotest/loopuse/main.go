package main

import (
    "fmt"
)

func main() {
    a := ClassA{}
    b := a.Generator()
    fmt.Println(a.id, b.id)
}
