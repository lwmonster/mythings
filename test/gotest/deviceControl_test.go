package test

import (
    "testing"
    "encoding/json"
    "fmt"
    "bytes"
    "net/http"
    "strconv"
)

func TestControlDevice(t *testing.T) {
    deviceId := "B4430DB16D230000"
    action := "on"
    controlURL := "http://test.ablecloud.cn:5000/DemoService/v2/sendToDevice"
    contentType := "application/x-zc-object"
    majorDomain := "ablecloud"
    subDomain := "test"
    accessMode := 1
    httpClient := &http.Client{}

    reqBody := make(map[string]string)
    reqBody["physicalDeviceId"] = deviceId
    reqBody["action"] = action
    b, err := json.Marshal(reqBody)
    if(err != nil) {
        fmt.Printf("error when encode reqBody, err:[%v]\n", err)
        return
    }
    reader := bytes.NewReader(b)
    req, err := http.NewRequest("POST", controlURL, reader)
    if err != nil {
        fmt.Printf("new request failed, err:[%v]\n", err)
        return
    }
    req.Header.Add("Content-Type", contentType)
    req.Header.Add("X-Zc-Major-Domain", majorDomain)
    req.Header.Add("X-Zc-Sub-Domain", subDomain)
    req.Header.Add("X-Zc-Access-Mode", strconv.Itoa(accessMode))


    resp, err := httpClient.Do(req)
    if err != nil {
        fmt.Printf("Do http request error, err:[%v]\n", err)
        return
    }
    fmt.Printf("response: [%v]", resp)

    if controlResult, ok := resp.Header["X-Zc-Msg-Name"]; ok {
        fmt.Println("control result:", controlResult)
    }
}

