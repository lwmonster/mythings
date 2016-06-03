package main

import (
    "fmt"
    "reflect"
)

func main() {
    var arr []interface{}

    arr = append(arr, 1)
    arr = append(arr, 2)
    arr = append(arr, 3)
    arr = append(arr, 4)
    arr = append(arr, "5")
    arr = append(arr, "6")
    arr = append(arr, "7")
    arr = append(arr, "8")

    fmt.Println(arr)


    for idx, v := range arr {
        //fmt.Printf("type of v: %T \n", v.(float64))
        fmt.Printf("%v value is %v\n", idx, v)
    }

    var a interface{}
    a = "ssss"

    switch t := a.(type) {
    case string:
        fmt.Println("arr is a string")
    case int:
        fmt.Println("arr is an integer")
    default:
        fmt.Printf("type is %v\n", t)
    }

    t := reflect.TypeOf(arr)
    fmt.Println(t)
    k := t.Kind()
    fmt.Println(k)

    x, ok := a.([]interface{})
    if !ok {
        fmt.Println(x)
        fmt.Println("a is not a slice:", a)
    } else {
        fmt.Println("good of a:", a)
    }



    y := 66.666

    y = int64(y)
    fmt.Println(y)
}
