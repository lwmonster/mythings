package main

import (
    "sync"
    "time"
    "math/rand"
    "fmt"
    "container/list"
    "math"
    "encoding/json"
    "net/http"
    "bytes"
    "io/ioutil"
)

type Location struct {
    Country     string  `json:"country"`
    Province    string  `json:"province"`
    City        string  `json:"city"`
}

/** 设备信息 **/
type DeviceProfile struct {
    Id       int64    `json:"id"`
    Loc      Location `json:"location"`
}

/** 用户信息 **/
type UserProfile struct {
    Id          int64      `json:"id"`
    Name        string     `json:"name"`
    Gender      string     `json:"gender"`
    Age         int        `json:"age"`
    Birthday    string     `json:"birthday"`
    Height      float64    `json:"height"`
    Weight      float64    `json:"weight"`
    Bmi         float64    `json:"bmi"`
    Fatratio    float64    `json:"fatratio"`
    Email       string     `json:"email"`
    Tel         string     `json:"tel"`
    Loc         Location   `json:"location"`
}


type App struct {
    uuid    string  `json:"uuid"`
    version string  `json:"version"`
    name    string  `json:"name"`
    os      string  `json:"os"`
    model   string  `json:"model"`
}

type  UserBasicInfo struct {
    name    string
    gender  string
}

type Event struct {
    TimeStamp string `json:"timestamp"`
}

type DeviceActiveEvent struct {
    Evt Event  `json:"event"`
    Did int64   `json:"did"`
}

type UserRegisterEvent struct {
    Evt Event `json:"event"`
    Uid int64 `json:"uid"`
}

type DeviceBindEvent struct {
    Evt Event `json:"event"`
    Did int64 `json:"did"`
    Uid int64 `json:"uid"`
    Uuid string `json:"uuid"`
    Status string `json:"status"`
}

type Mobile struct {
    Model string `json:"model"`
}

type AppSend struct {
    Name string `json:"name"`
    Version  string  `json:"version"`
    OS string `json:"os"`
}

type AppOperateEvent struct {
    Evt Event `json:"event"`
    MobilePhone Mobile `json:"mobile"`
    Application AppSend  `json:"app"`
    Op  string     `json:"op"`
    Uid int64       `json:"uid"`
    Uuid string     `json:"uuid"`
}

type AppAccessEvent struct {
    Evt Event `json:"event"`
    MobilePhone Mobile `json:"mobile"`
    Application AppSend  `json:"app"`
    Op  string     `json:"op"`
    Uid int64       `json:"uid"`
    Uuid string     `json:"uuid"`
    PageUrl string  `json:"pageurl"`
}

type MeasurementEvent struct {
    Evt Event `json:"event"`
    Uid int64 `json:"uid"`
    Did int64 `json:"did"`
    Weight float64 `json:"weight"`
    FatRatio float64 `json:"fatratio"`
}

type Operation struct {
    CMD string `json:"cmd"`
}
type DeviceOperateEvent struct {
    Evt Event `json:"event"`
    Uid int64 `json:"uid"`
    Did int64 `json:"did"`
    Op  Operation `json:"op"`
}

type ProfileActor struct {
}



var (
    appVersion  = []string{"1.0", "2.0", "2.0", "2.0", "3.0", "3.0", "3.0", "3.0", "3.0"}

    mobiles      = [][2]string{
        {"ios", "iPhone"},
        {"ios", "iPhone"},
        {"ios", "iPhone"},
        {"ios", "iPhone"},
        {"ios", "iPhone"},
        {"andriod", "Nexus 5"},
        {"andriod", "MX5"},
        {"andriod", "xiaomi"},
        {"andriod", "xiaomi"},
        {"andriod", "SUMSANG"},
        {"andriod", "SUMSANG"},
        {"andriod", "vivio"},
        {"andriod", "OPPO"},
        {"andriod", "SMARTISAN"},
        {"andriod", "HUAWEI"},
        {"andriod", "HUAWEI"},
        {"andriod", "Lenovo"},
        {"andriod", "COOLPAD"},
        {"andriod", "LETV"},
    }



    users       = []UserBasicInfo {
        {"Andrew", "male"},
        {"Andy", "male"},
        {"Arthur", "male"},
        {"Bob", "male"},
        {"Bill", "male"},
        {"Blake", "male"},
        {"Carl", "male"},
        {"Cary", "male"},
        {"Clark", "male"},
        {"Cole", "male"},
        {"Colin", "male"},
        {"Daniel", "male"},
        {"David", "male"},
        {"Duke", "male"},
        {"Esion", "male"},
        {"Evan", "male"},
        {"Frank", "male"},
        {"Gary", "male"},
        {"Glen", "male"},
        {"Harry", "male"},
        {"Tom", "male"},
        {"Jim", "male"},
        {"Jack", "male"},
        {"Landnerd", "male"},
        {"Cooker", "male"},
        {"Ray", "male"},
        {"James", "male"},
        {"Mark", "male"},
        {"Steve", "male"},
        {"Abby", "female"},
        {"Angel", "female"},
        {"Anna", "female"},
        {"Bette", "female"},
        {"Bonnie", "female"},
        {"Candy", "female"},
        {"Cherry", "female"},
        {"Chris", "female"},
        {"Daisy", "female"},
        {"Ellen", "female"},
        {"Emily", "female"},
        {"Emma", "female"},
        {"Eva", "female"},
        {"Fay", "female"},
        {"Grace", "female"},
        {"Iris", "female"},
        {"Hillary", "female"},
        {"Jade", "female"},
        {"Jane", "female"},
        {"Joan", "female"},
        {"Kitty", "female"},
        {"Mandy", "female"},
        {"May", "female"},
    }
)

var places = []Location{
    {"中国", "北京", "北京"},
    {"中国", "上海", "上海"},
    {"中国", "天津", "天津"},
    {"中国", "重庆", "重庆"},
    {"中国", "北京", "北京"},
    {"中国", "上海", "上海"},
    {"中国", "天津", "天津"},
    {"中国", "北京", "北京"},
    {"中国", "上海", "上海"},
    {"中国", "内蒙古", "呼和浩特"},
    {"中国", "内蒙古", "包头"},
    {"中国", "内蒙古", "赤峰"},
    {"中国", "内蒙古", "通辽"},
    {"中国", "内蒙古", "呼伦贝尔"},
    {"中国", "黑龙江", "哈尔滨"},
    {"中国", "黑龙江", "齐齐哈尔"},
    {"中国", "黑龙江", "牡丹江"},
    {"中国", "辽宁", "沈阳"},
    {"中国", "辽宁", "大连"},
    {"中国", "辽宁", "阜新"},
    {"中国", "吉林", "长春"},
    {"中国", "吉林", "吉林"},
    {"中国", "河北", "石家庄"},
    {"中国", "河北", "保定"},
    {"中国", "河南", "郑州"},
    {"中国", "河南", "南阳"},
    {"中国", "山东", "济南"},
    {"中国", "山东", "青岛"},
    {"中国", "山西", "太原"},
    {"中国", "山西", "大同"},
    {"中国", "湖北", "武汉"},
    {"中国", "湖北", "咸宁"},
    {"中国", "湖北", "黄石"},
    {"中国", "湖北", "荆门"},
    {"中国", "湖南", "长沙"},
    {"中国", "湖南", "常德"},
    {"中国", "四川", "成都"},
    {"中国", "江西", "南昌"},
    {"中国", "江西", "九江"},
    {"中国", "安徽", "合肥"},
    {"中国", "江苏", "南京"},
    {"中国", "江苏", "苏州"},
    {"中国", "浙江", "杭州"},
    {"中国", "浙江", "台州"},
    {"中国", "广东", "广州"},
    {"中国", "广东", "深圳"},
    {"中国", "广西", "南宁"},
    {"中国", "广西", "桂林"},
    {"中国", "云南", "昆明"},
    {"中国", "贵州", "贵阳"},
    {"中国", "西藏", "拉萨"},
    {"中国", "青海", "西宁"},
    {"中国", "甘肃", "兰州"},
    {"中国", "陕西", "西安"},
    {"中国", "宁夏", "银川"},
    {"中国", "新疆", "乌鲁木齐"},
}

func GetOneLocation() *Location {
    idx := rand.Intn(len(places))
    return &places[idx]
}

func Round(val float64, roundOn float64, places int ) (newVal float64) {
    var round float64
    pow := math.Pow(10, float64(places))
    digit := pow * val
    _, div := math.Modf(digit)
    if div >= roundOn {
        round = math.Ceil(digit)
    } else {
        round = math.Floor(digit)
    }
    newVal = round / pow
    return
}

//////////////////////////////////////////////////////////////////////
type Generator struct {
    deviceNum   int
    userNum     int
    appNum      int
    deviceList  *list.List
    userList    *list.List
    appList     *list.List

    deviceCur   *list.Element
    userCur     *list.Element
    appCur      *list.Element

    daysInterval int
    startTime   time.Time
    endTime     time.Time
    timeFormat  string
    zipfRegisterTime *rand.Zipf

    httpclient  *http.Client
}

func (gen *Generator) GenDeviceInfo() {
    for i := 0; i < gen.deviceNum; i++ {
        l := GetOneLocation()
        d := &DeviceProfile{Id: int64(i) + 200, Loc: *l}
        gen.deviceList.PushBack(d)
    }
    gen.deviceCur = gen.deviceList.Front()
}

func (gen *Generator) GetOneDevice() (*DeviceProfile, bool) {
    if gen.deviceCur != nil {
        tmp := gen.deviceCur.Value.(*DeviceProfile)
        gen.deviceCur = gen.deviceCur.Next()
        return tmp, true
    }
    return nil, false
}

func (gen *Generator) GetAge() int{
    //正态分布的年龄, 期望29, 标准差7
    return int(rand.NormFloat64() * 7 + 29)
}

func (gen *Generator) GetHeight(male bool) float64{
    //正态分布的身高,
    if male {
        return rand.NormFloat64() * 8 + 170
    }

    return rand.NormFloat64() * 8 + 156
}

func (gen *Generator) GetWeight(male bool) float64{
    //正态分布的体重
    if male {
        return rand.NormFloat64() * 8 + 66.2
    }

    return rand.NormFloat64() * 8 + 57.3
}

func (this *Generator) CalcBMIAndFatRatio(w, h, age float64, male bool) (float64, float64) {
    bmi := Round(w / (h*h/10000), .5, 2)
    fatratio := 1.2 * bmi + 0.23 * float64(age) - 5.4
    if male {
        fatratio = fatratio - 10.8
    }

    fatratio = Round(fatratio, .5, 2)

    return bmi, fatratio
}


func (gen *Generator) GenUserInfo() {
    now := time.Now()
    for i := 0; i < gen.userNum; i++ {
        idx := rand.Intn(len(users))
        u := &users[idx]
        loc := GetOneLocation()
        tmpUserProfile := &UserProfile{Id: int64(20500000 + i), Name: u.name, Gender: u.gender, Loc: *loc}
        male := true
        if tmpUserProfile.Gender == "female" {
            male = false
        }

        // 1. 生成age
        age := gen.GetAge()
        tmpBirthday := now.Add(-1 * time.Duration(age) * 365 * 24 * time.Hour)
        birthday := tmpBirthday.Format(gen.timeFormat)
        tmpUserProfile.Birthday = birthday
        tmpUserProfile.Age = age

        // 2. 生成身高
        height := gen.GetHeight(male)
        tmpUserProfile.Height = Round(height, .5, 1)

        // 3. 生成体重
        weight := gen.GetWeight(male)
        tmpUserProfile.Weight = Round(weight, .5, 1)

        // 4. 计算BMI
        bmi := weight / (height * height/10000)
        tmpUserProfile.Bmi = Round(bmi, .5, 2)

        // 5. 计算体脂率
        fatratio := 1.2 * bmi + 0.23 * float64(age) - 5.4
        if male {
            fatratio = fatratio - 10.8
        }
        tmpUserProfile.Fatratio = Round(fatratio, .5, 1)

        gen.userList.PushBack(tmpUserProfile)
    }

    gen.userCur = gen.userList.Front()
}

func (this *Generator) GetOneUser() (*UserProfile, bool) {
    if this.userCur != nil {
        tmp := this.userCur.Value.(*UserProfile)
        this.userCur = this.userCur.Next()
        return tmp, true
    }
    return nil, false
}

func (this *Generator) GenAppInfo() {
    for i := 0; i < this.appNum; i++ {
        appVersionIdx := rand.Intn(len(appVersion))
        mobileIdx := rand.Intn(len(mobiles))
        app := &App{
            uuid: fmt.Sprintf("f1a983l-aac4-462f-8426-%06d", i),
            version: appVersion[appVersionIdx],
            model: mobiles[mobileIdx][1],
            os:  mobiles[mobileIdx][0],
            name : "BodyFat",
        }

        this.appList.PushBack(app)
    }

    this.appCur = this.appList.Front()
}

func (this *Generator) GetOneApp() (*App, bool) {
    if this.appCur != nil {
        tmp := this.appCur.Value.(*App)
        this.appCur = this.appCur.Next()
        return tmp, true
    }
    return nil, false
}

func (this *Generator) PrintDeviceInfo() {
    for cur := this.deviceList.Front(); cur != nil; cur = cur.Next() {
        fmt.Println(cur.Value.(*DeviceProfile))
    }
}

func (this *Generator) PrintUserInfo() {
    for cur := this.userList.Front(); cur != nil; cur = cur.Next() {
        fmt.Println(cur.Value.(*UserProfile))
    }
}

func (this *Generator) PrintAPPInfo() {
    for cur := this.appList.Front(); cur != nil; cur = cur.Next() {
        fmt.Println(cur.Value.(*App))
    }
}

func (this *Generator) GetActiveTime() time.Time {
    // days * hour * minute * second
    // 秒
    interval := rand.Intn(this.daysInterval * 24 * 60 * 60)
    return this.startTime.Add(time.Duration(interval) * time.Second)
}

// 根据激活时间生成用户注册时间
func (this *Generator) GetRegisterTime(t time.Time) time.Time {
    // 间隔时间（分钟）
    interval := this.zipfRegisterTime.Uint64() + 10
    return t.Add(time.Duration(interval) * time.Minute)
}

func (this *Generator) GetBindTime(t time.Time) time.Time {
    interval := rand.Intn(13 * 60) // 秒
    return t.Add(time.Duration(interval) * time.Second)
}

func (this *Generator) GetUserBehaviorType() int {
    r := rand.Intn(100)
    if r < 25 {
        return 1    // 均匀型
    } else if r >= 25 && r < 50 {
        return 2    // 增加型
    } else {
        return 3    // 减弱型
    }
}


func (this *Generator) GenTimeSeries(t time.Time, userBehaviorType int) []int {
    var result []int

    timeInterval := this.endTime.Sub(t)
    daysInterval := int64(timeInterval.Hours()/24)

    //fmt.Println(timeInterval)
    //fmt.Println(daysInterval)

    // 使用频度
    strength := rand.Intn(100)
    // 绑定当天都要用一下
    result = append(result, 0)
    for i := 1; int64(i) < daysInterval; i++ {
        r := rand.Intn(100)
        if(r <= strength) {
            result = append(result, i)
        }
    }
    return result
}

func (this *Generator) GetUsingAppTimeLong() int64 {
    return int64(60 * (rand.NormFloat64() * 1.5 + 3))
}


func (this *Generator) isMeasure(t time.Time) bool {
    var prob int = 30
    hour := t.Hour()
    if (hour >= 17 && hour <= 24) {
        // 50% 的概率进行测量
        prob = 50
    }

    return rand.Intn(100) < prob
}


func (this *Generator) isShare(t time.Time) bool {
    var prob int = 10
    hour := t.Hour()
    if (hour >= 17 && hour <= 24) || (hour >= 0 && hour <=8) || (hour >= 11 && hour <= 13) {
        // 30% 的概率进行分享
        prob = 30
    }

    return rand.Intn(100) < prob
}


func (this *Generator) isLonin() bool {
    return rand.Intn(100) < 10
}


func (this *Generator) weightInc(u *UserProfile) bool {
    return rand.Intn(100) < 40
}

func (this *Generator) GetNewWeight(u *UserProfile, inc bool) float64 {
    if(inc) {
        return u.Weight + rand.Float64()
    }

    return u.Weight + rand.NormFloat64()*0.2 + 0.3
}

func (this *Generator) GenUserBehaviors(d *DeviceProfile, u *UserProfile, app *App, bdtime time.Time, tss []int) {
    weightInc := this.weightInc(u)

    for i := 0; i < len(tss); i++ {
        r := rand.Intn(24*60*60)
        startAppTime := bdtime.Add(time.Duration(tss[i]) * 24 *time.Hour).Add(time.Duration(r) * time.Second)
        tl := this.GetUsingAppTimeLong()
        endAppTime := startAppTime.Add(time.Duration(tl) * time.Second)

        fmt.Println(endAppTime)
        // 1. app_operate start
        aoes := &AppOperateEvent{
            Evt: Event{
                TimeStamp: startAppTime.Format(this.timeFormat),
            },

            MobilePhone: Mobile{
                Model: app.model,
            },

            Application: AppSend{
                Name: app.name,
                Version: app.version,
                OS: app.os,
            },

            Op: "start",
            Uid: u.Id,
            Uuid: app.uuid,
        }

        b, err := json.Marshal(aoes)
        if err != nil {
            panic("dump AppOperateEvent data error.")
        }
        fmt.Println(string(b))
        this.WriteEvent("app_operate", bytes.NewReader(b))

        // 2. app_operate end
        aoee := &AppOperateEvent{
            Evt: Event{
                TimeStamp: endAppTime.Format(this.timeFormat),
            },

            MobilePhone: Mobile{
                Model: app.model,
            },

            Application: AppSend{
                Name: app.name,
                Version: app.version,
                OS: app.os,
            },

            Op: "end",
            Uid: u.Id,
            Uuid: app.uuid,
        }

        b, err = json.Marshal(aoee)
        if err != nil {
            panic("dump AppOperateEvent data error.")
        }
        fmt.Println(string(b))
        this.WriteEvent("app_operate", bytes.NewReader(b))

        // 3. app_access  login
        loginTime := startAppTime.Add(time.Duration(rand.Intn(27) + 3) * time.Second)
        if this.isLonin() && loginTime.Before(endAppTime){
            aae := &AppAccessEvent{
                Evt: Event{
                    TimeStamp: loginTime.Format(this.timeFormat),
                },
                MobilePhone: Mobile{
                    Model: app.model,
                },

                Application: AppSend{
                    Name: app.name,
                    Version: app.version,
                    OS: app.os,
                },
                Op: "login",
                Uid: u.Id,
                Uuid: app.uuid,
                PageUrl: "main_page",
            }

            b, err = json.Marshal(aae)
            if err != nil {
                panic("dump AppAccessEvent data error.")
            }
            fmt.Println(string(b))
            this.WriteEvent("app_access", bytes.NewReader(b))
        }

        // 4. app_access  share
        shareTime := loginTime.Add(time.Duration(rand.Intn(100) + 10) * time.Second)
        if this.isShare(shareTime) && shareTime.Before(endAppTime) {
            aae := &AppAccessEvent{
                Evt: Event{
                    TimeStamp: loginTime.Format(this.timeFormat),
                },
                MobilePhone: Mobile{
                    Model: app.model,
                },

                Application: AppSend{
                    Name: app.name,
                    Version: app.version,
                    OS: app.os,
                },
                Op: "share",
                Uid: u.Id,
                Uuid: app.uuid,
                PageUrl: "measure_result",
            }

            b, err = json.Marshal(aae)
            if err != nil {
                panic("dump AppAccessEvent data error.")
            }
            fmt.Println(string(b))
            this.WriteEvent("app_access", bytes.NewReader(b))
        }
        // 5. measurement
        measureTime := loginTime.Add(time.Duration(rand.Intn(300) + 10) * time.Second)
        measure := this.isMeasure(measureTime)
        if measure && measureTime.Before(endAppTime) {
            male := true
            if u.Gender == "female" {
                male = false
            }
            userBirthday, err := time.Parse(this.timeFormat, u.Birthday)
            if err != nil {
                fmt.Println(err)
                panic("parse user birthday error")
            }
            userAge := time.Now().Year() - userBirthday.Year()
            currWeight := Round(this.GetNewWeight(u, weightInc), .5, 1)

            currBMI, currFatRatio := this.CalcBMIAndFatRatio(currWeight, u.Height, float64(userAge), male)

            u.Weight = currWeight
            u.Fatratio = currFatRatio
            u.Bmi = currBMI

            me := &MeasurementEvent{
                Evt : Event{
                    TimeStamp: measureTime.Format(this.timeFormat),
                },
                Uid : u.Id,
                Did : d.Id,
                Weight: currWeight,
                FatRatio: currFatRatio,
            }

            b, err = json.Marshal(me)
            if err != nil {
                panic("dump Measurement data error.")
            }
            fmt.Println(string(b))
            this.WriteEvent("measurement", bytes.NewReader(b))
        }

        // 6. device_operate start
        if measure {
            deviceStartTime := measureTime.Add(-1 * time.Duration(rand.Intn(200)) * time.Second)
            doe := &DeviceOperateEvent{
                Evt: Event{
                    TimeStamp: deviceStartTime.Format(this.timeFormat),
                },
                Uid : u.Id,
                Did : d.Id,

                Op : Operation{
                    CMD : "start",
                },
            }

            b, err = json.Marshal(doe)
            if err != nil {
                panic("dump Measurement data error.")
            }
            fmt.Println(string(b))
            this.WriteEvent("device_operate", bytes.NewReader(b))

        // 7. device_operate end
            deviceEndTime := deviceStartTime.Add(time.Duration(rand.Intn(300)) * time.Second)
            doe = &DeviceOperateEvent{
                Evt: Event{
                    TimeStamp: deviceEndTime.Format(this.timeFormat),
                },

                Uid : u.Id,
                Did : d.Id,

                Op : Operation{
                    CMD : "end",
                },
            }

            b, err = json.Marshal(doe)
            if err != nil {
                panic("dump Measurement data error.")
            }
            fmt.Println(string(b))
            this.WriteEvent("device_operate", bytes.NewReader(b))
        }


        //endAppTime.Hour()
    }
}

func (this *Generator) WriteEvent(collection string, reader *bytes.Reader) {
    //reqUrl := "http://10.254.178.55:10000/zc-queryengine/v1/write?collection=" + collection
    //reqUrl := "http://42.159.238.180:10000/zc-queryengine/v1/write?collection=" + collection
    reqUrl := "http://10.0.0.7:10000/zc-queryengine/v1/write?collection=" + collection

    r, err := http.NewRequest("POST", reqUrl, reader)
    if err != nil {
        panic("new request failed")
    }
    r.Header.Add("Content-Type", "text/json")
    r.Header.Add("X-Zc-Major-Domain-Id", "840")
    if resp, err := this.httpclient.Do(r); err != nil {
        panic(err)
    } else if resp.StatusCode != 200 {
        panic(resp.StatusCode)
    } else {
        _, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            panic(err)
        }
        resp.Body.Close()
    }
}

func (this *Generator) WriteProfile(writeType string, reader *bytes.Reader) {
    //reqUrl := "http://10.254.178.55:10000/zc-queryengine/v1/profile?type=" + writeType
    //reqUrl := "http://42.159.238.180:10000/zc-queryengine/v1/profile?type=" + writeType
    reqUrl := "http://10.0.0.7:10000/zc-queryengine/v1/profile?type=" + writeType

    r, err := http.NewRequest("POST", reqUrl, reader)
    if err != nil {
        panic("new request failed")
    }
    r.Header.Add("Content-Type", "text/json")
    r.Header.Add("X-Zc-Major-Domain-Id", "840")
    if resp, err := this.httpclient.Do(r); err != nil {
        panic(err)
    } else if resp.StatusCode != 200 {
        panic(resp.StatusCode)
    } else {
        _, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            panic(err)
        }
        resp.Body.Close()
    }
}


func (this *Generator) ProcessOneDevice() bool {
    //timeFormat := "2006-01-02 15:04:05.000000 -0700"
    d, ok := this.GetOneDevice()
    if !ok {
        //panic("No device available")
        return false
    }

    u, ok := this.GetOneUser()
    if !ok {
        //panic("No User available")
        return false
    }

    app, ok := this.GetOneApp()
    if !ok {
        //panic("No App available")
        return false
    }

    // 1. 激活设备
    activeTime := this.GetActiveTime()
    dae := &DeviceActiveEvent{
        Evt : Event{
            TimeStamp: activeTime.Format(this.timeFormat),
        },
        Did : d.Id,
    }
    b, err := json.Marshal(dae)
    if err != nil {
        panic("dump DeviceActiveEvent data error.")
    }

    fmt.Println(string(b))
    this.WriteEvent("device_activate", bytes.NewReader(b))
    // TODO  flush data to DB

    // 2. 注册用户
    registerTime := this.GetRegisterTime(activeTime)
    ure := &UserRegisterEvent{
        Evt : Event{
            TimeStamp: registerTime.Format(this.timeFormat),
        },
        Uid : u.Id,
    }

    b, err = json.Marshal(ure)
    if err != nil {
        panic("dump UserRegisterEvent data error.")
    }
    fmt.Println(string(b))

    this.WriteEvent("user_register", bytes.NewReader(b))
    // TODO  flush data to DB

    // 3. 绑定 设备、用户和app
    bindTime := this.GetBindTime(registerTime)
    dbe := &DeviceBindEvent{
        Evt : Event{
            TimeStamp: bindTime.Format(this.timeFormat),
        },
        Uid : u.Id,
        Did : d.Id,
        Uuid : app.uuid,
        Status : "bind",
    }

    b, err = json.Marshal(dbe)
    if err != nil {
        panic("dump DeviceBindEvent data error.")
    }
    fmt.Println(string(b))
    this.WriteEvent("device_bind", bytes.NewReader(b))

    // 4. profile_actor
    b, err = json.Marshal(u)
    if err != nil {
        panic("dump UserProfile data error.")
    }
    fmt.Println(string(b))
    this.WriteProfile("actor", bytes.NewReader(b))

    // 5. profile_object
    b, err = json.Marshal(d)
    if err != nil {
        panic("dump DeviceProfile data error.")
    }
    fmt.Println(string(b))
    this.WriteProfile("object", bytes.NewReader(b))


    // 6. 生成 用户的使用时间序列
    userUseTimeSeries := this.GenTimeSeries(bindTime, 1)
    fmt.Println(userUseTimeSeries)
    // 7. 生成用户行为
    this.GenUserBehaviors(d, u, app, bindTime, userUseTimeSeries)

    return true
}

func (this *Generator) Run() {
    for i := 0; i < this.deviceList.Len(); i++ {
        ok := this.ProcessOneDevice()
        if !ok {
            break
        }
    }
}


var generator *Generator
var once sync.Once

func GetGenerator() *Generator {
    once.Do(func (){
        src := rand.NewSource(time.Now().UnixNano())
        randGenerator := rand.New(src)

        v := 1.2
        s := 1.1
        st, _ := time.Parse("2006-01-02 15:04:05.000000 -0700", "2016-03-21 10:10:44.000000 +0800")
        days  := 65
        generator = &Generator{
            deviceNum   : 50000,
            userNum     : 50000,
            appNum      : 50000,
            deviceList  : list.New(),
            userList    : list.New(),
            appList     : list.New(),
            deviceCur   : nil,
            userCur     : nil,
            appCur      : nil,

            daysInterval: days,
            startTime   : st,
            endTime     : st.Add(time.Duration(days) * 24 * time.Hour),
            timeFormat  : "2006-01-02 15:04:05.000000 -0700",
            zipfRegisterTime: rand.NewZipf(randGenerator, v, s, 40),
            httpclient  : &http.Client{},
        }

        // init dataset
        generator.GenDeviceInfo()
        generator.GenUserInfo()
        generator.GenAppInfo()
    })

    fmt.Println(generator)
    return generator
}


//////////////////////////////////////////////////////////////////////
func main() {
    rand.Seed(time.Now().UnixNano())
    mygenerator := GetGenerator()

    //mygenerator.PrintDeviceInfo()
    //mygenerator.PrintUserInfo()
    //mygenerator.PrintAPPInfo()

    /**
    for i := 0; i < 50; i++ {
        activeTime := mygenerator.GetActiveTime()
        registerTime := mygenerator.GetRegisterTime(activeTime)
        bindTime    :=  mygenerator.GetBindTime(registerTime)
        fmt.Println("active time:", activeTime.String())
        fmt.Println("register time:", registerTime.String())
        fmt.Println("bind time:", bindTime.String())
    }
    **/
    //mygenerator.ProcessOneDevice()

    mygenerator.Run()
}

