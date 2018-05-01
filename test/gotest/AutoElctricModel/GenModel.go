package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
    "time"
    "strconv"
    "math"
    qecommon "ac-bigdata/ac-queryengine/common"
    bigcommon "ac-bigdata-common/common"
    "flag"
    "net/http"
    "bytes"
    "encoding/json"
)

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

    return round / pow
}

const (
    NoAction string = "noaction"
    NeedTurnOn string = "on"
    NeedTurnOff string = "off"
)


type DayTime struct {
    Hour   int
    Minute int
}

func (this *DayTime) Before(dt DayTime) bool {
    if(this.Hour < dt.Hour) {
        return true
    }

    if(this.Hour > dt.Hour) {
        return false
    }

    if(this.Minute < dt.Minute) {
        return true
    }

    return false
}

func (this *DayTime) After(dt DayTime) bool {
    return dt.Before(*this)
}


func (this *DayTime) Equal(dt DayTime) bool {
    return this.Hour == dt.Hour && this.Minute == dt.Minute
}

type ModelValue struct {
    Start       DayTime
    End         DayTime
    ElecMount   float64
    Days        int
    AvgElecMount    float64

    tmpElecMount    float64
    tmpRecordCount  int
}

type ElectricController struct {
    DataFileName    string
    ModelDir        string

    DBConfigFileName string

    TimeFormat      string
    TimeZone        string
    StartTime       time.Time
    TimeInterval    time.Duration

    RecordInterval  int  // 上报间隔 单位(秒)
    RecordsNumDur   int  // 一个统计区间中的上报条数

    NoUseStartDayTime DayTime
    NoUseEndDayTime   DayTime


    controlURL  string
    majorDomain string
    subDomain   string
    accessMode  int
    contentType string
    deviceId    string

    httpClient  *http.Client
    // 处理Model的区块索引
    processIdx      int
    Threshold       float64
    Model           []ModelValue
}


func (this *ElectricController) Init(du time.Duration, confFile string) {
    this.DataFileName = "/home/liangwei/gowork/fakedata"
    this.ModelDir = "/home/ablecloud/bin/predictor/model/"
    //this.DBConfigFileName = "/home/ablecloud/bin/queryengine/conf"
    //this.DBConfigFileName = "/home/liangwei/gowork/src/ac-bigdata/ac-queryengine/conf/queryengine.conf"
    this.DBConfigFileName = confFile
    config, err := qecommon.LoadConfig(this.DBConfigFileName)
    if(err != nil) {
        fmt.Printf("load config error err:[%v]\n", err)
        return
    }
    fmt.Println("Load config sucess")

    bigcommon.InitDBHandler(config.VerticaConfig)

    this.TimeFormat = "2006-01-02 15:04:05.000000-0700"
    this.TimeZone = "Asia/Shanghai"
    //this.TimeZone = "Asia/Tokyo"
    this.StartTime, _ = time.Parse(this.TimeFormat, "2000-01-01 00:00:00.000000+0800")
    // 区间 30 分钟
    this.TimeInterval = du
    // 设备上报间隔
    this.RecordInterval = 3
    // 一个区间中的记录条数
    this.RecordsNumDur = int( this.TimeInterval / (time.Duration(this.RecordInterval) * time.Second) )

    // 不使用的时间段， 据此调整threshold
    this.NoUseStartDayTime = DayTime{
        Hour: 0,
        Minute: 0,
    }
    this.NoUseEndDayTime = DayTime{
        Hour: 9,
        Minute: 30,
    }

    this.controlURL = "http://test.ablecloud.cn:5000/DemoService/v2/sendToDevice"
    this.majorDomain = "ablecloud"
    this.subDomain = "test"
    this.accessMode = 1
    this.contentType = "application/x-zc-object"
    this.deviceId = "B4430DB16D230000"
    this.httpClient = &http.Client{}


    this.processIdx = 0
    this.Threshold = 0.0
    this.Model = make([]ModelValue, 0)
    tmpTime := time.Date(this.StartTime.Year(), this.StartTime.Month(), this.StartTime.Day(), 0, 0, 0, 0, this.StartTime.Location() )
    tmpEndTime := tmpTime.Add(time.Duration(1) * 24 * time.Hour).Add(time.Duration(-1) * time.Second)
    for currTime := tmpTime; currTime.Before(tmpEndTime); {
        currEndTime := currTime.Add(this.TimeInterval)
        endHour := currEndTime.Hour()
        if(currTime.Hour() > endHour) {
            endHour = 24
        }
        mv := &ModelValue{
            Start: DayTime{
                Hour: currTime.Hour(),
                Minute: currTime.Minute(),
            },
            End: DayTime{
                Hour: endHour,
                Minute: currEndTime.Minute(),
            },
            ElecMount: float64(0.0),
            Days: 0,
            AvgElecMount: float64(0.0),
            tmpElecMount: float64(0.0),
            tmpRecordCount: 0,
        }

        this.Model = append(this.Model, *mv)

        currTime = currEndTime
    }
}


func (this *ElectricController) GetDataFromDBnTrainModel() {
    sqlStr := `SELECT "event.timestamp", total_quantity, use_quantity `
    sqlStr += `FROM "3_plug_electric" `
    sqlStr += `WHERE "event.timestamp" > '` + this.StartTime.Format(this.TimeFormat) + `' `
    sqlStr += `ORDER BY "event.timestamp" limit 50000`

    tx, err := bigcommon.GetDBHandler().Begin()
    if err != nil {
        fmt.Printf("Set timezone failed, err:[%v]\n", err)
        return
    }
    defer tx.Commit()
    _, err = tx.Exec(fmt.Sprintf("SET TIMEZONE TO '%v'", this.TimeZone))
    if err != nil {
        fmt.Printf("Set timezone failed, err:[%v]\n", err)
        return
    }

    //rows, err := bigcommon.GetDBHandler().Query(sqlStr)
    rows, err := tx.Query(sqlStr)
    if(err != nil) {
        fmt.Printf("SELECT FROM DB error err:[%v]\n", err)
        return
    }
    defer rows.Close()

    //_, offset := this.StartTime.Zone()
    //fmt.Println("offset:", offset)
    for rows.Next() {
        var timeValue time.Time
        var totalElecMount float64
        var useElecMount   float64

        err := rows.Scan(&timeValue, &totalElecMount, &useElecMount)
        if err != nil {
            fmt.Printf("Fetch result failed, err:[%v]\n", err)
            return
        }

        //fmt.Println(timeValue.Format(this.TimeFormat))
        //fmt.Println("select location:", timeValue.Location())

        //_, offset2 := timeValue.Zone()
        //fmt.Println("select offset2:", offset2)
        //timeValue = timeValue.Add(time.Duration(-offset) * 1e9)
        //timeValue = timeValue.Add(time.Duration(offset2) * 1e9)
        //timeValue = timeValue.In(this.StartTime.Location())

        if(timeValue.Sub(this.StartTime) >= this.TimeInterval ||
                timeValue.Weekday() == time.Sunday) {
            this.StartTime = timeValue
            continue
        }

        this.StartTime = timeValue

        this.UpdateModel(timeValue, totalElecMount, useElecMount)
        fmt.Println(timeValue.Format(this.TimeFormat), totalElecMount, useElecMount)
    }
}


func (this *ElectricController) GetElectricData(startTime time.Time) {

    if file, err := os.Open(this.DataFileName); err == nil {
        defer file.Close()
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            cols := strings.Split(scanner.Text(), "\t")
            if(len(cols) != 3) {
                continue
            }

            currTime, err := time.Parse(this.TimeFormat, cols[0])
            if(err != nil) {
                fmt.Printf("parse time error, timeStr:[%s] err:[%v]", cols[0], err)
            }
            // 更新开始时间，
            this.StartTime = currTime
            useElecMount, err := strconv.ParseFloat(cols[2], 64)
            if(err != nil) {
                fmt.Printf("parse electric mount error, electricMount:[%s] err:[%v]", cols[2], err)
                continue
            }
            totalElecMount, err := strconv.ParseFloat(cols[1], 64)
            if(err != nil) {
                fmt.Printf("parse electric mount error, electricMount:[%s] err:[%v]", cols[1], err)
                continue
            }

            this.UpdateModel(currTime, totalElecMount, useElecMount)
        }
    }
}

func (this *ElectricController) UpdateModel(currTime time.Time, totalElecMount, useElecMount float64) {
    tmpDayTime := DayTime{
        Hour: currTime.Hour(),
        Minute: currTime.Minute(),
    }

    for {
        // 寻找区间块
        tmpModelValue := &this.Model[this.processIdx]
        if tmpModelValue.Start.Equal(tmpDayTime) ||
                (tmpModelValue.Start.Before(tmpDayTime) &&
                tmpModelValue.End.After(tmpDayTime)) {
            tmpModelValue.tmpElecMount += useElecMount
            tmpModelValue.tmpRecordCount += 1
            break
        }

        // 走到这里说明要换区间块， 需要把当前区间的ModelValue更新
        if(tmpModelValue.tmpRecordCount > 0 &&
                tmpModelValue.tmpRecordCount > int(float64(this.RecordsNumDur) * 0.8)) {

            tmpModelValue.ElecMount += tmpModelValue.tmpElecMount
            tmpModelValue.Days += 1
            tmpModelValue.AvgElecMount = tmpModelValue.ElecMount / float64(tmpModelValue.Days)
        }

        if(tmpModelValue.tmpRecordCount <= int(float64(this.RecordsNumDur) * 0.8)) {
            fmt.Printf("discard data for record num insufficient, date:%d-%02d-%02d timeRegion:[%02d:%02d ~%02d:%02d] elecMount:%.4f record num:%d\n",
                    currTime.Year(), currTime.Month(), currTime.Day(),
                    tmpModelValue.Start.Hour, tmpModelValue.Start.Minute, tmpModelValue.End.Hour, tmpModelValue.End.Minute,
                    tmpModelValue.tmpElecMount, tmpModelValue.tmpRecordCount)
        }

        // 不管有没有更新，都重置 ModelValue 的tmp变量
        // 这里的tmp变量主要是为了过滤出现异常数据不足的情况
        tmpModelValue.tmpRecordCount = 0
        tmpModelValue.tmpElecMount = float64(0.0)

        this.processIdx = (this.processIdx + 1) % len(this.Model)
    }

}


func (this *ElectricController) UpdateThreshold() {
    ElecMount := float64(0.0)
    regionNum := 0
    /**
    for _, mv := range this.Model {
        if(mv.Start.Equal(this.NoUseStartDayTime) || mv.Start.After(this.NoUseStartDayTime)) &&
            mv.End.Before(this.NoUseEndDayTime) {
            ElecMount += mv.ElecMount
            regionNum += mv.Days
        }
    }
    **/


    for _, mv := range this.Model {
        ElecMount += mv.ElecMount
        regionNum += mv.Days
    }


    if(regionNum > 0) {
        this.Threshold = ElecMount / float64(regionNum)
    }
}

func (this *ElectricController) OutputModel() {

    t := time.Now()
    filenameSufix := fmt.Sprintf("%04d-%02d-%02d.%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    if file, err := os.Create(this.ModelDir + "/model." + filenameSufix); err == nil {
        defer file.Close()
        file.WriteString(fmt.Sprintf("StartTime:%s\n", this.StartTime.Format(this.TimeFormat)))
        file.WriteString(fmt.Sprintf("Threshold:%.4f\n", this.Threshold))
        for _, mv := range this.Model {
            file.WriteString(fmt.Sprintf("%d\t%d\t%d\t%d\t%.4f\t%d\t%.4f\t%.4f\t%d\n",
                mv.Start.Hour, mv.Start.Minute, mv.End.Hour, mv.End.Minute,
                mv.ElecMount, mv.Days, mv.AvgElecMount, mv.tmpElecMount, mv.tmpRecordCount))
            //fmt.Println(mv.Start, mv.End, Round(mv.ElecMount, .5, 10), mv.Days, Round(mv.AvgElecMount, .5, 10), mv.tmpElecMount, mv.tmpRecordCount)
        }
    } else {
        fmt.Printf("Create file error err:%v\n", err)
    }
}

func (this *ElectricController) CheckStatus(tm time.Time) string {
    dtm := &DayTime{
        Hour: tm.Hour(),
        Minute: tm.Minute(),
    }

    tmpDts := DayTime{
        Hour: 22,
        Minute: 0,
    }
    tmpDte := DayTime{
        Hour: 24,
        Minute: 0,
    }

    if(dtm.Equal(tmpDts) || dtm.After(tmpDts)) && dtm.Before(tmpDte) {
        return NeedTurnOff
    }

    if(dtm.Equal(this.NoUseStartDayTime) || dtm.After(this.NoUseStartDayTime)) &&
            dtm.Before(this.NoUseEndDayTime) {
        return NeedTurnOff
    }

    for i := 0; i < len(this.Model); i++ {
        tmpModelValue := &this.Model[i]
        if dtm.Equal(tmpModelValue.Start) ||
                (dtm.After(tmpModelValue.Start) && (dtm.Before(tmpModelValue.End))) {
            if(tmpModelValue.Days == 0) {
                return NoAction
            }

            if(tmpModelValue.AvgElecMount > this.Threshold) {
                return NeedTurnOn
            }

            return NeedTurnOff
        }
    }

    return NoAction
}

func (this *ElectricController) doControl(action string) {
    if(action != NeedTurnOn && action != NeedTurnOff) {
        fmt.Printf("%s, do nothing to device:%s\n", time.Now().Format(this.TimeFormat), this.deviceId)
        return
    }

    reqBody := make(map[string]string)
    reqBody["physicalDeviceId"] = this.deviceId
    reqBody["action"] = action
    b, err := json.Marshal(reqBody)
    if (err != nil) {
        fmt.Printf("error when encode reqBody, err:[%v]\n", err)
        return
    }
    reader := bytes.NewReader(b)
    req, err := http.NewRequest("POST", this.controlURL, reader)
    if err != nil {
        fmt.Printf("new request failed, err:[%v]\n", err)
        return
    }
    req.Header.Add("Content-Type", this.contentType)
    req.Header.Add("X-Zc-Major-Domain", this.majorDomain)
    req.Header.Add("X-Zc-Sub-Domain", this.subDomain)
    req.Header.Add("X-Zc-Access-Mode", strconv.Itoa(this.accessMode))

    resp, err := this.httpClient.Do(req)
    if err != nil {
        fmt.Printf("Do http request error, err:[%v]\n", err)
        return
    }
    fmt.Printf("%s turn [%s] device:[%s] sucess response: [%v]\n", time.Now().Format(this.TimeFormat), action, this.deviceId, resp)
}


func (this *ElectricController) ControlOn() {
    for {
        tm := time.Now()
        action := this.CheckStatus(tm)

        if(action == NeedTurnOn) {
            this.doControl(action)
        }

        // 一分钟检查一次
        time.Sleep(time.Duration(5) * time.Second)
    }
}


func (this *ElectricController) ControlOff() {
    now := time.Now()
    action := this.CheckStatus(now)
    if(action == NeedTurnOff) {
        this.doControl(action)
    } else {
        fmt.Printf("%s ControlOff do nothing\n", now.Format(this.TimeFormat))
    }

    triggerTime := now.Truncate(time.Hour).Add(this.TimeInterval).Add(time.Duration(1)*time.Minute)
    fmt.Println(triggerTime)

    if triggerTime.Before(now) {
        triggerTime = triggerTime.Add(this.TimeInterval)
    }
    time.Sleep(triggerTime.Sub(now))

    for {
        tm := time.Now()
        action := this.CheckStatus(tm)
        if(action == NeedTurnOff) {
            this.doControl(action)
        } else {
            fmt.Printf("%s ControlOff do nothing\n", tm.Format(this.TimeFormat))
        }

        // 半小时 检查一次
        time.Sleep(this.TimeInterval)
    }
}

func (this *ElectricController) Run() {
    go this.ControlOn()
    go this.ControlOff()
    for {
        this.GetDataFromDBnTrainModel()
        this.UpdateThreshold()
        this.OutputModel()

        time.Sleep(time.Duration(60) * time.Second)
    }
}




func main() {

    confFileName := flag.String("conf", "/home/ablecloud/bin/queryengine/conf/queryengine.conf", "config file")
    //var confFileName = flag.String("conf", "/home/liangwei/gowork/src/ac-bigdata/ac-queryengine/conf/queryengine.conf", "config file")
    controller := &ElectricController{}
    controller.Init(time.Duration(30) * time.Minute, *confFileName)

    controller.Run()
}

