package main

import (
    "container/list"
    "time"
    "net/http"
    "math/rand"
    "fmt"
    "sync"
    "encoding/json"
    "bytes"
    "io/ioutil"
    "strconv"

    "ac-common-go/developer"
)

type Location struct {
    Country     string  `json:"country"`
    Province    string  `json:"province"`
    City        string  `json:"city"`
}

/** 设备信息 **/
type DeviceProfile struct {
    Id       int64    `json:"id"`
    Model    string   `json:"model"`
    Loc      Location `json:"location"`
}

type  UserBasicInfo struct {
    name    string
    gender  string
}

/** 用户信息 **/
type UserProfile struct {
    Id          int64      `json:"id"`
    Name        string     `json:"name"`
    Gender      string     `json:"gender"`
    Birthday    string     `json:"birthday"`
    Age         int        `json:"age"`
    Email       string     `json:"email"`
    Tel         string     `json:"tel"`
    Loc         Location   `json:"location"`
}

type PhoneBasicInfo struct {
    Company     string     `json:"company"`
    Model       string     `json:"model"`
    OS          string     `json:"os"`
    OSVersion   string     `json:"osversion"`
}

/** 手机信息 **/
type PhoneProfile struct {
    Id          string     `json:"id"`
    Company     string     `json:"company"`
    Model       string     `json:"model"`
    OS          string     `json:"os"`
    OSVersion   string     `json:"osversion"`
    Loc         Location   `json:"location"`
}

/** app信息 **/
type App struct {
    Name        string     `json:"name"`
    Version     string     `json:"version"`
}

/****/
type Event struct {
    TimeStamp string `json:"timestamp"`
}

type DeviceActiveEvent struct {
    Evt Event  `json:"event"`
    Did int64  `json:"did"`
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

type Operation struct {
    CMD string `json:"cmd"`
}
type DeviceOperateEvent struct {
    Evt Event `json:"event"`
    Uid int64 `json:"uid"`
    Did int64 `json:"did"`
    Op  Operation `json:"op"`
}

type CookEvent struct {
    Evt Event `json:"event"`
    Uid int64 `json:"uid"`
    Did int64 `json:"did"`
    Op  string `json:"op"`
    Mode string `json:"mode"`
}

type RecipeEvent struct {
    Evt Event `json:"event"`
    Uid int64       `json:"uid"`
    Uuid string     `json:"uuid"`
    Op   string     `json:"op"`
    RecipeType string     `json:"recipe_type"`
    RecipeURL  string     `json:"recipe_url"`
}

type DeviceFaultEvent struct {
    Evt Event `json:"event"`
    Did int64 `json:"did"`
    FaultCode int `json:"fault_code"`
    CookStatus string `json:"cook_status"`
    CookMode string `json:"cook_mode"`
}

type AppOperateEvent struct {
    Evt Event `json:"event"`
    Op  string     `json:"op"`
    Uid int64       `json:"uid"`
    Uuid string     `json:"uuid"`
    MobileModel string `json:"mobile.model"`
    AppName string  `json:"app.name"`
    AppOS   string  `json:"app.os"`
    AppVersion string  `json:"app.version"`
}

type AppAccessEvent struct {
    Evt Event `json:"event"`
    Op  string     `json:"op"`
    Uid int64       `json:"uid"`
    Uuid string     `json:"uuid"`
    PageUrl string  `json:"pageurl"`
}

var cookMode = []string{
    "rice",
    "rice_pudding",
    "gruel",
    "soup",
    "steam",
}

var recipeTypes = []string {
    "chuan", "xiang", "yue", "lu", "dongbei", "xibei", "zhe",
}


var (
    appVersion  = []string{"1.0", "2.0", "2.0", "2.0", "3.0", "3.0", "3.0", "3.0", "3.0"}

    mobiles      = []PhoneBasicInfo{
        {"apple", "iPhone", "ios", "9.1"},
        {"apple", "iPhone", "ios", "9.1"},
        {"apple", "iPhone", "ios", "9.1"},
        {"apple", "iPhone", "ios", "9.1"},
        {"apple", "iPhone", "ios", "9.1"},
        {"apple", "iPhone", "ios", "9.1"},
        {"apple", "iPhone", "ios", "9.1"},
        {"Google", "Nexus 5", "andriod", "5.5"},
        {"MeiZu", "MX5", "andriod", "5.1"},
        {"xiaomi", "MI 4", "andriod", "5.5"},
        {"xiaomi", "HM note", "andriod", "5.5"},
        {"SUMSANG", "Galarx6", "andriod", "5.5"},
        {"SUMSANG", "Galarx7", "andriod", "5.5"},
        {"bubugao", "vivio", "andriod", "5.5"},
        {"bubugao", "OPPO", "andriod", "5.5"},
        {"chuizi", "SMARTISAN", "andriod", "5.5"},
        {"HUAWEI", "Mate7", "andriod", "4.5"},
        {"HUAWEI", "Mate8", "andriod", "5.5"},
        {"Lenovo", "Lenovo3", "andriod", "5.5"},
    }



    users = []UserBasicInfo {
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
    {"中国", "北京", "北京"},
    {"中国", "上海", "上海"},
    {"中国", "北京", "北京"},
    {"中国", "上海", "上海"},
    {"中国", "北京", "北京"},
    {"中国", "上海", "上海"},
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
    {"中国", "内蒙古", "通辽"},
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
    {"中国", "广东", "广州"},
    {"中国", "广东", "广州"},
    {"中国", "广东", "广州"},
    {"中国", "广东", "广州"},
    {"中国", "广东", "深圳"},
    {"中国", "广东", "深圳"},
    {"中国", "广东", "深圳"},
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

func GetOneUser() *UserBasicInfo{
    idx := rand.Intn(len(users))
    return &users[idx]
}

func GetOnePhone() *PhoneBasicInfo{
    idx := rand.Intn(len(mobiles))
    return &mobiles[idx]
}

//////////////////////////////////////////////////////////////////////
type Generator struct {
    deviceNum   int
    userNum     int
    phoneCurrId int
    deviceList  *list.List
    userList    *list.List
    //appList     *list.List

    deviceCur   *list.Element
    userCur     *list.Element

    daysInterval int
    startTime   time.Time
    endTime     time.Time
    timeFormat  string
    zipfRegisterTime *rand.Zipf

    httpclient  *http.Client

    majorDomain string
    subDomain   string
    developerId int64
    ak          string
    sk          string
    timeOut     int
}

func (this *Generator) WriteEvent(collection string, reader *bytes.Reader) {
    //reqUrl := "http://10.254.178.55:10000/zc-queryengine/v1/write?collection=" + collection
    timestamp := time.Now().Unix()
    method := "write"
    nonce := "abc123"

    signer := &developer.Signer{
        DeveloperId : int64(this.developerId),
        MajorDomain : this.majorDomain,
        SubDomain   : this.subDomain,
        Timestamp   : timestamp,
        Timeout     : int64(this.timeOut),
        Nonce       : nonce,
        Method      : method,
    }
    signString, err := signer.Sign(this.sk)
    /**
    fmt.Println("DeveloperId:", this.developerId)
    fmt.Println("MajorDomain:", this.majorDomain)
    fmt.Println("SubDomain:", this.subDomain)
    fmt.Println("Timestamp:", timestamp)
    fmt.Println("Timeout:", this.timeOut)
    fmt.Println("Nonce:", nonce)
    fmt.Println("Method:", method)
    fmt.Println("signString:", signString)
    **/
    if(err != nil) {
        fmt.Println("gen signature failed, %v", err)
        panic("gen signature failed")
    }

    reqUrl := "http://test.ablecloud.cn:5000/zc-queryengine/v1/" + method + "?collection=" + collection

    r, err := http.NewRequest("POST", reqUrl, reader)
    if err != nil {
        panic("new request failed")
    }
    r.Header.Add("Content-Type", "text/json")
    r.Header.Add("X-Zc-Major-Domain", this.majorDomain)
    r.Header.Add("X-Zc-Sub-Domain", this.subDomain)
    r.Header.Add("X-Zc-Developer-Id", strconv.Itoa(int(this.developerId)))
    r.Header.Add("X-Zc-Access-Key", this.ak)
    r.Header.Add("X-Zc-Timestamp", strconv.FormatInt(timestamp, 10))
    r.Header.Add("X-Zc-Timeout", strconv.Itoa(this.timeOut))
    r.Header.Add("X-Zc-Nonce", nonce)
    r.Header.Add("X-Zc-Developer-Signature", signString)


    if resp, err := this.httpclient.Do(r); err != nil {
        panic(err)
    } else if resp.StatusCode != 200 {
        panic(resp.StatusCode)
    } else {
        result, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            panic(err)
        }
        resp.Body.Close()
        fmt.Println(string(result))
    }
}

func (this *Generator) WriteProfile(writeType string, reader *bytes.Reader) {
    timestamp := time.Now().Unix()
    method := "profile"
    nonce := "abc123"

    signer := &developer.Signer{
        DeveloperId : int64(this.developerId),
        MajorDomain : this.majorDomain,
        SubDomain   : this.subDomain,
        Timestamp   : timestamp,
        Timeout     : int64(this.timeOut),
        Nonce       : nonce,
        Method      : method,
    }
    signString, err := signer.Sign(this.sk)

    fmt.Println("DeveloperId:", this.developerId)
    fmt.Println("MajorDomain:", this.majorDomain)
    fmt.Println("SubDomain:", this.subDomain)
    fmt.Println("Timestamp:", timestamp)
    fmt.Println("Timeout:", this.timeOut)
    fmt.Println("Nonce:", nonce)
    fmt.Println("Method:", method)
    fmt.Println("signString:", signString)

    if(err != nil) {
        fmt.Println("gen signature failed, %v", err)
        panic("gen signature failed")
    }

    //reqUrl := "http://10.254.178.55:10000/zc-queryengine/v1/" + method + "?type=" + writeType
    reqUrl := "http://test.ablecloud.cn:5000/zc-queryengine/v1/" + method + "?type=" + writeType

    r, err := http.NewRequest("POST", reqUrl, reader)
    if err != nil {
        panic("new request failed")
    }
    r.Header.Add("Content-Type", "text/json")
    r.Header.Add("X-Zc-Major-Domain", this.majorDomain)
    r.Header.Add("X-Zc-Sub-Domain", this.subDomain)
    r.Header.Add("X-Zc-Developer-Id", strconv.Itoa(int(this.developerId)))
    r.Header.Add("X-Zc-Access-Key", this.ak)
    r.Header.Add("X-Zc-Timestamp", strconv.FormatInt(timestamp, 10))
    r.Header.Add("X-Zc-Timeout", strconv.Itoa(this.timeOut))
    r.Header.Add("X-Zc-Nonce", nonce)
    r.Header.Add("X-Zc-Developer-Signature", signString)
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


func (this *Generator) GetAge(isMale bool) int{
    if isMale {
        //正态分布的年龄, 期望30, 标准差10
        return int(rand.NormFloat64() * 10 + 30)
    }

    return int(rand.NormFloat64() * 12 + 27)
}

func (this *Generator) GenDeviceInfo() {

    for i := 0; i < this.deviceNum; i++ {
        //l := GetOneLocation()
        d := &DeviceProfile{Id: int64(i) + 100200, Model:"E8"}
        this.deviceList.PushBack(d)
    }
    this.deviceCur = this.deviceList.Front()
}


func (this *Generator) GenUserInfo() {
    now := time.Now()
    for i := 0; i < this.userNum; i++ {
        u := GetOneUser()
        //loc := GetOneLocation()
        tmpUserProfile := &UserProfile{Id: int64(20100000 + i), Name: u.name, Gender: u.gender}
        male := true
        if tmpUserProfile.Gender == "female" {
            male = false
        }

        // 1. 生成age
        age := this.GetAge(male)
        tmpBirthday := now.Add(-1 * time.Duration(age) * 365 * 24 * time.Hour)
        birthday := tmpBirthday.Format(this.timeFormat)
        tmpUserProfile.Birthday = birthday
        tmpUserProfile.Age = age

        this.userList.PushBack(tmpUserProfile)
    }

    this.userCur = this.userList.Front()
}


func (this *Generator) GetOneDevice() (*DeviceProfile, bool) {
    if this.deviceCur != nil {
        tmp := this.deviceCur.Value.(*DeviceProfile)
        this.deviceCur = this.deviceCur.Next()
        return tmp, true
    }
    return nil, false
}

func (this *Generator) GetUser() (*UserProfile, bool) {
    if this.userCur != nil {
        tmp := this.userCur.Value.(*UserProfile)
        this.userCur = this.userCur.Next()
        return tmp, true
    }
    return nil, false
}

func (this *Generator) GetPhone() *PhoneProfile {
    id := fmt.Sprintf("f1a983l-aac4-462f-8426-%06d", this.phoneCurrId)
    this.phoneCurrId += 1
    pbi := GetOnePhone()
    return &PhoneProfile{Id:id, Company:pbi.Company, Model: pbi.Model, OS: pbi.OS, OSVersion: pbi.OSVersion}
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

func (this *Generator) GenTimeSeries(t time.Time) []int {
    var result []int

    timeInterval := this.endTime.Sub(t)
    daysInterval := int64(timeInterval.Hours()/24)

    //fmt.Println(timeInterval)
    //fmt.Println(daysInterval)

    // 使用频度
    //strength := rand.Intn(100)
    weekday_strength := rand.NormFloat64()*5 + 30
    weekend_strength := rand.NormFloat64()*10 + 60
    // 绑定当天都要用一下
    result = append(result, 0)
    for i := 1; int64(i) < daysInterval; i++ {
        r := rand.Intn(100)
        currTime := t.Add(time.Duration(i) * 24 * time.Hour)
        if(currTime.Weekday() == 0 || currTime.Weekday() == 1) {
            if(float64(r) <= weekend_strength) {
                result = append(result, i)
            }
        } else {
            if(float64(r) <= weekday_strength) {
                result = append(result, i)
            }
        }
    }
    return result
}


func (this *Generator) GetUsingAppTimeLong() int64 {
    return int64(60 * (rand.NormFloat64() * 1.5 + 3))
}


func (this *Generator) WriteDeviceOperateEvt(d *DeviceProfile,
        u *UserProfile, p *PhoneProfile, operateTime time.Time, cmd string) {
    doe := &DeviceOperateEvent{
        Evt: Event{
            TimeStamp: operateTime.Format(this.timeFormat),
        },
        Uid : u.Id,
        Did : d.Id,

        Op : Operation{
            CMD : cmd,
        },
    }

    b, err := json.Marshal(doe)
    if err != nil {
        panic("dump device operate event error.")
    }
    fmt.Println(string(b))
    this.WriteEvent("device_operate", bytes.NewReader(b))
}


func (this *Generator) WriteCookEvt(d *DeviceProfile,
        u *UserProfile, p *PhoneProfile, operateTime time.Time, ckmd string , cmd string) {

    e := &CookEvent{
        Evt: Event{
            TimeStamp: operateTime.Format(this.timeFormat),
        },
        Uid: u.Id,
        Did: d.Id,
        Op: cmd,
        Mode: ckmd,
    }

    b, err := json.Marshal(e)
    if err != nil {
        panic("dump cook event error.")
    }
    fmt.Println(string(b))
    this.WriteEvent("device_cook", bytes.NewReader(b))
}


func (this *Generator) WriteDeviceFaultEvt(d *DeviceProfile,
        u *UserProfile, p *PhoneProfile, operateTime time.Time, tmpCookMode string,
        reserving, cooking, warming bool) {

    commonErrorCode := []int {2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, }
    reserveErrorCode := []int {7001, 7002, 7003, 7004, }
    reserveErrorCode = append(reserveErrorCode, commonErrorCode...)
    cookErrorCode := []int {8001, 8002, 8003, 8004, }
    cookErrorCode = append(cookErrorCode, commonErrorCode...)
    warmErrorCode := []int {9001, 9002, 9003, 9004, }
    warmErrorCode = append(warmErrorCode, commonErrorCode...)

    ecode := 0
    cookStatus  := ""
    if(reserving) {
        ecode = reserveErrorCode[rand.Intn(len(reserveErrorCode))]
        cookStatus = "reserving"
    }
    if(cooking) {
        ecode = cookErrorCode[rand.Intn(len(cookErrorCode))]
        cookStatus = "cooking"
    }
    if(warming) {
        ecode = warmErrorCode[rand.Intn(len(warmErrorCode))]
        cookStatus = "warming"
    }

    e := &DeviceFaultEvent{
        Evt: Event{
            TimeStamp: operateTime.Format(this.timeFormat),
        },

        Did: d.Id,
        FaultCode: ecode,
        CookStatus: cookStatus,
        CookMode: tmpCookMode,
    }

    b, err := json.Marshal(e)
    if err != nil {
        panic("dump device fault event error.")
    }
    fmt.Println(string(b))
    this.WriteEvent("device_fault", bytes.NewReader(b))
}

func (this *Generator) WriteAppOperateEvt(d *DeviceProfile,
    u *UserProfile, p *PhoneProfile, operateTime time.Time, op string) {

    e := &AppOperateEvent{
        Evt: Event{
            TimeStamp: operateTime.Format(this.timeFormat),
        },
        Uid: u.Id,
        Uuid: p.Id,
        Op: op,
        MobileModel: p.Model,
        AppName: "cooker",
        AppOS: p.OS,
        AppVersion: "3.2",
    }

    b, err := json.Marshal(e)
    if err != nil {
        panic("dump app operate err.")
    }
    fmt.Println(string(b))
    this.WriteEvent("app_operate", bytes.NewReader(b))
}

func (this *Generator) WriteAppAccessEvt(d *DeviceProfile,
        u *UserProfile, p *PhoneProfile, operateTime time.Time, op, page string) {

    e := &AppAccessEvent{
        Evt: Event{
            TimeStamp: operateTime.Format(this.timeFormat),
        },
        Uid: u.Id,
        Uuid: p.Id,
        Op: op,
        PageUrl: page,
    }

    b, err := json.Marshal(e)
    if err != nil {
        panic("dump app access err.")
    }
    fmt.Println(string(b))
    this.WriteEvent("app_access", bytes.NewReader(b))
}

func (this *Generator) WriteRecipeEvt(d *DeviceProfile,
        u *UserProfile, p *PhoneProfile, operateTime time.Time, op, recipeType, recipeURL string) {

    e := &RecipeEvent{
        Evt: Event{
            TimeStamp: operateTime.Format(this.timeFormat),
        },
        Uid: u.Id,
        Uuid: p.Id,
        Op: op,
        RecipeType: recipeType,
        RecipeURL: recipeURL,
    }

    b, err := json.Marshal(e)
    if err != nil {
        panic("dump recipe event err.")
    }
    fmt.Println(string(b))
    this.WriteEvent("recipe", bytes.NewReader(b))

}
func (this *Generator) GenOnceBehavior(d *DeviceProfile,
        u *UserProfile, p *PhoneProfile, operateTime time.Time, launch, dinner, night bool) {
    opTime := operateTime
    if(launch) {
        opTime = time.Date(operateTime.Year(), operateTime.Month(),
                operateTime.Day(), 11, 0, 0, 0, operateTime.Location())
    } else if (dinner) {
        opTime = time.Date(operateTime.Year(), operateTime.Month(),
            operateTime.Day(), 18, 0, 0, 0, operateTime.Location())

    } else if (night) {
        opTime = time.Date(operateTime.Year(), operateTime.Month(),
            operateTime.Day(), 22, 0, 0, 0, operateTime.Location())
    }

    // 1. device_operate start
    this.WriteDeviceOperateEvt(d, u, p, opTime, "start")


    // 2. device_cook  reserve
    ckmd := cookMode[rand.Intn(len(cookMode))]
    reserved := false
    reserveTime := opTime
    cookStartTime := opTime
    r := rand.Intn(100)
    if((launch || dinner ) && r <= 10) || (night && r <= 60) {
        tmpTimeLong := rand.Intn(120)
        reserveTime = opTime.Add(time.Duration(tmpTimeLong) * time.Second)
        this.WriteCookEvt(d, u, p, reserveTime, ckmd, "start_reserve")
        reserved = true
    }

    // 3. device_fault  1/1000 的概率发生错误
    if(reserved && rand.Intn(1000) == 1) {
        faultTime := reserveTime.Add(time.Duration(3) * 60 * time.Second)
        this.WriteDeviceFaultEvt(d, u, p, faultTime, ckmd, true, false, false)
        return
    }

    // 4. device_cook reserve end  && cook start
    if(reserved ) {
        rd := 0
        if(launch || dinner) {
            rd = rand.Intn(100) + 30
        } else {
            rd = rand.Intn(10) + 5
            if(rd > 10) {
                rd = 10
            }
            rd = rd * 60
        }
        reserveEndTime := reserveTime.Add(time.Duration(rd) * 60 * time.Second)
        cookStartTime = reserveEndTime
        this.WriteCookEvt(d, u, p, reserveEndTime, ckmd, "end_reserve")
        // start cook
        this.WriteCookEvt(d, u, p, cookStartTime, ckmd, "start_cook")
    } else  {
        cookStartTime = opTime.Add(time.Duration(3) * 60 * time.Second)
        this.WriteCookEvt(d, u, p, cookStartTime, ckmd, "start_cook")
    }


    // 5. cooking time fault
    if(rand.Intn(10000) == 1) {
        faultTime := cookStartTime.Add(time.Duration(rand.Intn(50)) * 60 * time.Second)
        this.WriteDeviceFaultEvt(d, u, p, faultTime, ckmd, false, true, false)
        return
    }


    // 6. cook end
    cookEndTime := cookStartTime.Add(time.Duration(rand.Intn(30) + 30) * 60 * time.Second)
    this.WriteCookEvt(d, u, p, cookEndTime, ckmd, "end_cook")

    // 7. warming start
    warmingStartTime := cookEndTime
    this.WriteCookEvt(d, u, p, warmingStartTime, ckmd, "start_warming")

    // 8. warming time fault
    if(rand.Intn(10000) == 1) {
        faultTime := warmingStartTime.Add(time.Duration(rand.Intn(10)) * 60 * time.Second)
        this.WriteDeviceFaultEvt(d, u, p, faultTime, ckmd, false, false, true)
        return
    }

    // 9. warming end
    warmingEndTime := warmingStartTime.Add(time.Duration(rand.Intn(20)) * 60 * time.Second)
    this.WriteCookEvt(d, u, p, warmingEndTime, ckmd, "end_warming")

    // 10. device_operate  end
    deviceOperateEndTime := warmingEndTime.Add(time.Duration(rand.Intn(100)) * 60 * time.Second)
    this.WriteDeviceOperateEvt(d, u, p, deviceOperateEndTime , "end")

    if(rand.Intn(100) <= 50) {

        // 5. app_operate
        appOperateTime := opTime.Add(time.Duration(rand.Intn(5)) * 60 * time.Second)
        this.WriteAppOperateEvt(d, u, p, appOperateTime, "start")

        // 6. app_access
        if (rand.Intn(100) < 1) {
            // login
            appOperateTime = appOperateTime.Add(time.Duration(rand.Intn(10)) * time.Second)
            this.WriteAppAccessEvt(d, u, p, appOperateTime, "login", "main_page")
        }

        if (rand.Intn(100) < 1) {
            appOperateTime = appOperateTime.Add(time.Duration(rand.Intn(20)) * 60 * time.Second)
            this.WriteAppAccessEvt(d, u, p, appOperateTime, "share", "recipe")
        }

        // 7. recipe
        if (rand.Intn(100) < 10) {
            appOperateTime = appOperateTime.Add(time.Duration(rand.Intn(20)) * 60 * time.Second)
            recipeType := recipeTypes[rand.Intn(len(recipeTypes))]
            recipeURL := recipeType + strconv.Itoa(rand.Intn(99))
            this.WriteRecipeEvt(d, u, p, appOperateTime, "brower", recipeType, recipeURL)

            if (rand.Intn(100) < 10) {
                appOperateTime = appOperateTime.Add(time.Duration(2) * 60 * time.Second)
                this.WriteRecipeEvt(d, u, p, appOperateTime, "collect", recipeType, recipeURL)
            }
        }

        // 8. end app
        appOperateTime = appOperateTime.Add(time.Duration(rand.Intn(3)) * 60 * time.Second)
        this.WriteAppOperateEvt(d, u, p, appOperateTime, "end")
    }
}


func (this *Generator) GenUserBehavior(device *DeviceProfile,
        user *UserProfile, phone *PhoneProfile, bdtime time.Time) {

    tss := this.GenTimeSeries(bdtime)

    for i := 0; i < len(tss); i++ {
        operateTime := bdtime.Add(time.Duration(tss[i]) * 24 * time.Hour)

        launch := false
        dinner := false
        nightSupply := false

        r := rand.Intn(100)
        if(r <= 20) {
            launch = true
        }

        r = rand.Intn(100)
        if(r <= 70) {
            dinner = true
        }

        r = rand.Intn(100)
        if(r <= 30) {
            nightSupply = true
        }

        if(launch) {
            this.GenOnceBehavior(device, user, phone, operateTime, launch, false, false)
        }

        if(dinner) {
            this.GenOnceBehavior(device, user, phone, operateTime, false, dinner, false)
        }

        if(nightSupply) {
            this.GenOnceBehavior(device, user, phone, operateTime, false, false, nightSupply)
        }

    }

}


func (this *Generator) ProcessOneDevice() bool {
    loc := GetOneLocation()
    d, ok := this.GetOneDevice()
    if !ok {
        fmt.Println("no device available")
        return false
    }

    u, ok := this.GetUser()
    if !ok {
        fmt.Println("no user available")
        return false
    }

    phone := this.GetPhone()

    d.Loc = *loc
    u.Loc = *loc
    phone.Loc = *loc

    fmt.Println("testing..........")
    fmt.Println("device:", *d)
    fmt.Println("user:", *u)
    fmt.Println("phone:", *phone)
    fmt.Println("testing..........")

    // 1. 设备激活
    activateTime := this.GetActiveTime()
    dae := &DeviceActiveEvent{
        Evt : Event{
            TimeStamp: activateTime.Format(this.timeFormat),
        },
        Did : d.Id,
    }
    b, err := json.Marshal(dae)
    if err != nil {
        panic("dump DeviceActiveEvent data error.")
    }
    fmt.Println(string(b))
    this.WriteEvent("device_activate", bytes.NewReader(b))

    // 2. 用户注册
    registerTime := this.GetRegisterTime(activateTime)
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

    // 3. 绑定
    bindTime := this.GetBindTime(registerTime)
    dbe := &DeviceBindEvent{
        Evt : Event{
            TimeStamp: bindTime.Format(this.timeFormat),
        },
        Uid : u.Id,
        Did : d.Id,
        Uuid : phone.Id,
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


    this.GenUserBehavior(d, u, phone, bindTime)


    return true
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
        st, _ := time.Parse("2006-01-02 15:04:05.000000 -0700", "2016-04-05 10:10:44.000000 +0800")
        days  := 95
        generator = &Generator{
            deviceNum   : 200000,
            userNum     : 200000,
            phoneCurrId : 100100,
            deviceList  : list.New(),
            userList    : list.New(),
            deviceCur   : nil,
            userCur     : nil,

            daysInterval: days,
            startTime   : st,
            endTime     : st.Add(time.Duration(days) * 24 * time.Hour),
            timeFormat  : "2006-01-02 15:04:05.000000 -0700",
            zipfRegisterTime: rand.NewZipf(randGenerator, v, s, 40),
            httpclient  : &http.Client{},

            majorDomain : "liangwei",
            subDomain   : "",
            developerId : 348,
            ak          : "3845acab407f9eac80d4294d424637a9",
            sk          : "abaadede40c5431980a8a4ff8edbf204",
            timeOut     : 3,
        }

        // init dataset
        generator.GenDeviceInfo()
        generator.GenUserInfo()
    })

    fmt.Println(generator)
    return generator
}


func main() {
    mygenerator := GetGenerator()


    mygenerator.Run()
}

