package test

import (
    "testing"
    "time"
    "fmt"
)

func TestTime(t *testing.T) {
    fmt.Println("come in ....")
    now := time.Now()

    fmt.Println(now)
    fmt.Println(now.Second())



    testTime, err := time.Parse("2006-01-02 15:04:05.000000-0700", "2016-05-07 10:14:54.03652+08")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(testTime)


    for i := 0; i < 10; i++ {
        currTime := now.Add(time.Duration(i) * 24 * time.Hour)

        fmt.Println(int(currTime.Weekday()))
        fmt.Println(currTime.Year(), int(currTime.Month()), currTime.Day(), currTime.Hour(), currTime.Minute(), currTime.Second())

        tmpTime := time.Date(currTime.Year(), currTime.Month(), currTime.Day(), 0, 0, 0, 0, currTime.Location())
        fmt.Println(tmpTime)
    }



    //
    fmt.Println("======================================================")

    now = time.Now()
    triggerTime := now.Truncate(time.Hour).Add(time.Duration(31)*time.Minute)
    fmt.Println(triggerTime)

    if triggerTime.Before(now) {
        triggerTime = triggerTime.Add(time.Duration(30) * time.Minute)
    }

    fmt.Println("trigger time:", triggerTime, "sleep:", triggerTime.Sub(now))

}
