package main

import (
    "math/rand"
    "time"
)

type Checker struct {
}


func (this *Checker) Check() *ClassA{
    seed := time.Now().Unix()
    rand.Seed(seed)
    id := rand.Intn(100000)
    return &ClassA{
        id: id,
    }
}
