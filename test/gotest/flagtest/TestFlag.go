package main

import (
    "flag"
    "fmt"
)

func main() {
    var x = flag.String("conf_file", "aaa", "config file")

    fmt.Println(*x)
}
