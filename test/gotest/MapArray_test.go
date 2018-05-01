package test

import (
    "testing"
    "fmt"
)

func TestMapArray(t *testing.T) {
    c := make([]map[string]interface{}, 0)

    v := make(map[string]interface{})
    c = append(c, v)

    v["aaa"] = 44

    if _, ok := v["bbb"] ; !ok {
        fmt.Println("no bbb in v")
        //v["bbb"] = make([]string, 1)
        v["bbb"] = 99
    }

    fmt.Println(c)


    for key, value := range v{
        fmt.Println(key,":", value)
    }
}
