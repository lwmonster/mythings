package main

import (
    "container/list"
    "math/rand"
    "time"
    "fmt"
)

var (
    lst *list.List
    currEle *list.Element
)

type MyStruct struct {
    id  int
    name string
}

func initList() {
    lst = list.New()
    for i := 0; i < 10; i++ {
        ms := &MyStruct{
            id : i,
            name : "aaaa",
        }
        lst.PushBack(ms)
    }

    currEle = lst.Front()
}


func getOneElement() (*MyStruct, bool){
    if(currEle != nil){
        //v := int(currEle.Value)
        v := currEle.Value.(*MyStruct)
        currEle = currEle.Next()
        return v, true
    }

    return nil, false
}


func main() {
    rand.Seed(time.Now().UnixNano())
    initList()



    //测试 游标顺序取
    for {
        r := rand.Intn(100)
        if r < 100 {
            x, ok := getOneElement()
            if(ok) {
                fmt.Println("get from list:", x.id, x.name)
                x.id = x.id*10
                x.name = x.name + "_test"
            } else {
                fmt.Println("no data available")
                break
            }
        } else {
            fmt.Println("don't get ")
        }
    }

    for e := lst.Front(); e != nil; e = e.Next() {
        fmt.Println(e.Value.(*MyStruct))
    }
}


