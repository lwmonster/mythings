package test

import (
    "fmt"
    "testing"
)

func TestConvertToString(t *testing.T){

    var v interface{}
    //v = make([]int, 0)
    //v = append(v, 100)
    v = 100
    //fmt.Println(v.(string))

    s, ok := v.(string)
    if(!ok) {
        t.Errorf("error")
    } else {
        fmt.Println(s)
    }
}
