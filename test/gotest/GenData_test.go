package test

import (
    "testing"
    "time"
    "fmt"
    "math/rand"
)


func getElectricMount(tm time.Time) float64 {
    sortTm := tm.Hour() * 100 + tm.Minute()
    if((sortTm >= 0 && sortTm <= 940) ||
            (sortTm > 1210 && sortTm <= 1300) ||
            (sortTm >= 1800 && sortTm <= 1850) ||
            (sortTm >= 2110 && sortTm <= 2359)) {
        prob := rand.Intn(100)
        if(prob < 80) {
            return float64(0.0)
        }
        return float64(0.01)
    }

    prob := rand.Intn(100)
    if(prob < 50) {
        return float64(0.0)
    }
    return float64(0.01)
}

func TestGenData(t *testing.T) {
    timeFormat := "2006-01-02 15:04:05.000000-0700"
    startTimeStr := "2016-05-02 00:00:01.000000+0800"

    endTimeStr := "2016-05-06 00:00:00.000000+0800"

    startTime, _ := time.Parse(timeFormat, startTimeStr)
    endTime, _ := time.Parse(timeFormat, endTimeStr)

    sum := float64(0.0)
    for currTime := startTime; currTime.Before(endTime); currTime = currTime.Add(time.Duration(3) * time.Second){
        elcMount := getElectricMount(currTime)
        sum += elcMount
        fmt.Printf("%s\t%v\t%v\n", currTime.Format(timeFormat), sum, elcMount)
    }
}