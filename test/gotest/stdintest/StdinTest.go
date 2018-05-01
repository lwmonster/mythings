package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
    "time"
    "strconv"
)

type TimeControl struct {
    Hour   int
    Minute int
}

func main() {
    /**  Scanln usage
    for {
        str := ""
        _, err := fmt.Scanln(&str)
        if(err != nil) {
            fmt.Println(err)
            break
        }
        fmt.Println(str)
    }
    **/

    model := make(map[[2]time.Time]float64)

    var startTime time.Time = nil
    var endTime time.Time = nil
    timeInterval  := time.Duration(30) * 60 *time.Second
    elcMount := float64(0.0)
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        //fmt.Println(scanner.Text())
        cols := strings.Split(scanner.Text(), "\t")
        if(len(cols) != 3) {
            continue
        }

        currTime, err := time.Parse("2006-01-02 15:04:05.000000 -0700", cols[0])
        if(err != nil) {
            fmt.Printf("parse time err, timeStr:[%s] err:[%v]\n", cols[0], err)
            continue
        }

        if(startTime == nil || currTime.Sub(startTime) > timeInterval) {
            startTime = currTime
            endTime = startTime.Add(timeInterval)
        }

        mt, err := strconv.ParseFloat(cols[1], 64)
        if(err != nil) {
            fmt.Printf("parse time err, timeStr:[%s] err:[%v]\n", cols[0], err)
            continue
        }
        elcMount += strconv.ParseFloat(cols[1], 64)



        fmt.Println(len(cols))
    }
}

