package main

import (
    "fmt"
    "time"
)

func main() {
    now := time.Now()
    fmt.Println("now:", now)
    fmt.Println("now utc:", now.UTC())
    fmt.Println("20 days later:", now.Add(20 * 24 * time.Hour))
    fmt.Println("20 days ago:", now.Add(-20 * 24 * time.Hour))

    fmt.Println("timestamp:", now.Unix())
    fmt.Println(now.Clock())
    fmt.Println("year:", now.Year())
    fmt.Println("month:", now.Month())
    fmt.Println("day:", now.Day())
    fmt.Println(now.Date())

    fmt.Println(now.Format("2006-01-01"))

    t, err := time.Parse("2006-01-02", "1987-09-28")
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(t)
    }


    t, err = time.Parse("2006-01-02 15:04:05.00000", "1984-11-07 10:11:33.12345")
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(t)
    }

    fmt.Println("now:", now.Round(time.Hour * 24))
    fmt.Println("now:", now.Truncate(time.Hour * 24))


    timeInterval := time.Now().Sub(t)
    fmt.Println(timeInterval.Hours())
}
