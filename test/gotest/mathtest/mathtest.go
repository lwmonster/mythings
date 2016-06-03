package main

import (
    "math/rand"
    "time"
    "fmt"
    "math"
)

/**
func  Round(x float64, n int) float64 {
    mod := math.Pow(10, float64(n))
    intVal := int64(x * mod)
    result := float64(intVal) / mod

    return result
}
**/

func Round(val float64, roundOn float64, places int ) float64 {
    var round float64
    pow := math.Pow(10, float64(places))
    digit := pow * val
    _, div := math.Modf(digit)
    //fmt.Println("div:", div)
    if div >= roundOn {
        round = math.Ceil(digit)
    } else {
        round = math.Floor(digit)
    }

    //fmt.Println("round:", round)

    return round / pow
}

func main() {
    src := rand.NewSource(time.Now().UnixNano())
    randGenerator := rand.New(src)

    /**
    //exp := 0.17
    for i := 0; i < 50; i++ {
        //fmt.Println(rand.ExpFloat64()/exp)
        fmt.Println(randGenerator.ExpFloat64())
    }
    **/

    fmt.Println("zipf distribution begin.......")
    v := 1.2
    s := 1.1
    zipfDistribution := rand.NewZipf(randGenerator, v, s, 40)

    for i := 0; i < 50; i++ {
        age := zipfDistribution.Uint64() + 10
        fmt.Println(age)
    }
    fmt.Println("zipf distribution end .......")


    fmt.Println("normal distribution end .......")

    for i := 0; i < 50; i++ {
        fmt.Println(Round(math.Abs(randGenerator.NormFloat64()*1.5 + 3) * 60, .5, 0))
    }


    fmt.Println("normal distribution end .......")

    fmt.Println("float round begin .......")
    f1 := 1.22499
    f2 := 1.22599
    f3 := 1.22699

    fmt.Printf("%v round:%v\n", f1, Round(f1, .5, 2))
    fmt.Printf("%v round:%v\n", f2, Round(f2, .5, 2))
    fmt.Printf("%v round:%v\n", f3, Round(f3, .5, 2))
    fmt.Println("float round end .......")
}


