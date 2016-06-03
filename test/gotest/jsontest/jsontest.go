package main

import (
    "fmt"
    "encoding/json"
)

type Event struct {
    TimeStamp  string `json:"timestamp"`
}

type DeviceActiveEvent struct {
    Event Event `json:"event"`
    deviceId string `json:"did"`
    测试字段 string
}

type Server struct {
    ServerName string
    ServerIP   string
}

type Serverslice struct {
    Servers []Server
}

func main() {
    var s Serverslice
    s.Servers = append(s.Servers, Server{ServerName: "Shanghai_VPN", ServerIP: "127.0.0.1"})
    s.Servers = append(s.Servers, Server{ServerName: "Beijing_VPN", ServerIP: "127.0.0.2"})
    b, err := json.Marshal(s)
    if err != nil {
        fmt.Println("json err:", err)
    }
    fmt.Println(string(b))

    dae := &DeviceActiveEvent{
        Event: Event{
            TimeStamp: "2016-01-05",
        },
        deviceId: "A99910",
        测试字段 : "测试",
    }

    fmt.Println(dae)
    daeStr, err := json.Marshal(dae)
    fmt.Println(string(daeStr))

}

/**
func main() {

    ts :=  time.Now().Format("2006-01-02 15:04:05.000000 +0800")

    de := &DeviceActiveEvent{
        event: Event{
            TimeStamp: ts,
        },
        deviceId: "10009",
    }


    fmt.Println(de)
    marshaled, err := json.Marshal(de)

    if err == nil {
        fmt.Println(string(marshaled))
    } else {
        fmt.Println(err)
    }

}

**/
