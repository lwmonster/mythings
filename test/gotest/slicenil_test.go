package test

import (
    "testing"
    "fmt"
)

func TestSliceNil(t *testing.T) {
    var slice []string

    slice = nil
    fmt.Println(len(slice))
    fmt.Println(slice == nil)

    for idx, e := range slice {
        fmt.Println(idx, e)
    }

    slice = append(slice, "aaa")
    fmt.Println(slice)


    var arr []int
    arr = make([]int, 5)
    arr = append(arr, 99)
    fmt.Println(arr)

    /**
    var slice1 = make([]string, 0)
    fmt.Println(slice1)
    fmt.Println(len(slice1))
    fmt.Println(slice1 == nil)

    var str string
    fmt.Println(str)
    fmt.Println(len(str))
    fmt.Println(str == nil)
    str = nil
    fmt.Println(str == nil)
    **/

}