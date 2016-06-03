package main

/**
    include :
        reading from stdin
        convert string to int
        string trim
        channel
 */

import (
    "fmt"
    /*
    "bufio"
    "os"
    "strconv"
    "strings"
    */
)

func put(c chan <- int) {
    for i := 0; i < 10; i++ {
        c <- i
    }
    close(c)
    /**
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("Enter text: ")
        text, _ := reader.ReadString('\n')
        text = strings.Trim(text, " \t\n")
        if text == "q" {
            close(c)
            break
        }

        fmt.Println("you input:", text)
        int_val, err := strconv.Atoi(text)
        if(err == nil) {
            c <- int_val
        }
    }
    **/
}

func  getAndPrint(c <- chan int) {
    /**
    for v := range(c){
        fmt.Println(v)
    }
    **/

    for {
        select {
        case v, ok:= <- c:
            if !ok {
                fmt.Println("will return")
                return
            }
            fmt.Println("read from channel:", v)
        }
    }
}

func main() {
    ch := make(chan int)
    put(ch)

    //getAndPrint(ch)
}


