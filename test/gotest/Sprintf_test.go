package test

import (
    "testing"
    "fmt"
    "time"
)

func genStr () string {
    return `abc %s %v`
}

func TestSprintf(t *testing.T) {
    s := genStr()

    s1 := fmt.Sprintf(s, "aa", 2)
    fmt.Println(s1)


    tm := time.Now()
    s2 := fmt.Sprintf("%04d-%02d-%02d.%02d:%02d:%02d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
    fmt.Println(s2)
}