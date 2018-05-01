package main

//import "golang.org/x/tour/tree"
import "tree"
import "fmt"

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
    if(t == nil) { return }

    Walk(t.Left, ch)

    ch <- t.Value

    Walk(t.Right, ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
    var ch1 = make(chan int)
    var ch2 = make(chan int)

    go Walk(t1, ch1)
    go Walk(t2, ch2)

    for i := 0; i < 10; i++ {
        v1 := <- ch1
        v2 := <- ch2
        fmt.Println("v1:", v1, "v2:", v2)
        if(v1 != v2) { return false }
    }
    return true
}

func main() {
    var t1 = tree.New(1)
    var t2 = tree.New(1)

    var t3 = tree.New(2)
    var t4 = tree.New(3)
    fmt.Println("t1:", t1)
    fmt.Println("t2:", t2)
    fmt.Println("t3:", t3)
    fmt.Println("t4:", t4)

    fmt.Println("t1 equals t2?", Same(t1, t2))
    fmt.Println("t3 equals t4?", Same(t3, t4))
}


