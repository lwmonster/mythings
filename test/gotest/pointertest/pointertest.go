package main

import (
    "container/list"
    "fmt"
)

type  Item struct {
    id int
    name string
}

var  (
   itemLst *list.List
)

func initList() {
    itemLst = list.New()
    /** 存储 指针
    itemLst.PushBack(&Item{1, "aaa"})
    itemLst.PushBack(&Item{2, "bbb"})
    itemLst.PushBack(&Item{3, "ccc"})
    **/

    // 不存指针
    itemLst.PushBack(Item{1, "aaa"})
    itemLst.PushBack(Item{2, "bbb"})
    itemLst.PushBack(Item{3, "ccc"})
}

func printList() {
    for e := itemLst.Front(); e != nil; e = e.Next() {
        v := e.Value.(Item)
        fmt.Println(v.id, v.name)
    }
}

func main(){
    initList()
    printList()
    v1 := itemLst.Front().Value.(Item)
    v1.id = 999
    printList()

    itemLst.Front().Value.(Item).id = 999
    printList()
}
