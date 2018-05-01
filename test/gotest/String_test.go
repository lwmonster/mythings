package test
import (
    "fmt"
    "testing"
    "time"
    "strconv"
)

func TestStringOperation112(t *testing.T) {
    s1 := `aaa`
    s2 := `bbb`

    s3 := "ccc " + s1 + " " + s2

    s4 := fmt.Sprintf(`4444 %s`, s1)
    fmt.Println(s3)
    fmt.Println(s4)

    now := time.Now()
    fmt.Println(now.Unix())
    //s := strconv.ParseInt(now.Unix(), 10, )
    s := strconv.FormatInt(now.Unix(), 10)
    fmt.Println(s)


    fmt.Printf("%.2f", 1.2282332)
}

