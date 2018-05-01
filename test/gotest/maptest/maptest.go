package main

import "fmt"

type DataMap map[string]interface{}

func testRef(dm DataMap) {
    dm["ccc"] = 999
}

func testPointer(dm *DataMap) {
    (*dm)["ddd"] = 888
}

func main() {
    dm := make(DataMap)

    dm["aaa"] = 100
    dm["bbb"] = "abc"


    fmt.Println(dm)

    testRef(dm)

    fmt.Println(dm)

    testPointer(&dm)

    fmt.Println(dm)
}
