package test

import (
    "testing"
    "fmt"
)

func TestFallThrough(t *testing.T) {
    a := 3
    switch a {
    case 2:
        fmt.Println("a = 2")
        fallthrough
    case 3:
        fmt.Println("a = 3")
        fallthrough
    case 4:
        fmt.Println("a = 4")
        fallthrough
    case 5:
        fmt.Println("a = 5")
    case 6:
        fmt.Println("a = 6")
    default:
        fmt.Println("a = default")

    }
}
