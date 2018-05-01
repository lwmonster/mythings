package main


type ClassA struct {
    id int
}

func (this *ClassA) Generator() *ClassA{
    checker := Checker{}
    return checker.Check()
}
