// time相关工具
package utils

import (
    "time"
    "fmt"
)

//json时间格式化
type JsonTime struct {
    time.Time
}

func NowTime() JsonTime {
    return JsonTime{time.Now()}
}


//实现json序列化方法,格式化时间
func (this JsonTime) MarshalJSON() ([]byte, error) {
    var stamp = fmt.Sprintf("\"%s\"", GetFormatYmdHisByTime(this.Time))
    return []byte(stamp), nil
}

// 获得YY-mm-dd时间格式
func GetFormatYmd() string {
    return GetFormatYmdByTime(time.Now())
}

// 根据Time获得YY-mm-dd时间格式
func GetFormatYmdByTime(t time.Time) string {
    return t.Format("2006-01-02")
}

// 根据时间戳获得YY-mm-dd时间格式
func GetFormatYmdByUnix(u int64) string {
    return GetFormatYmdByTime(time.Unix(u, 0))
}

// 获得YY-mm-dd HH:ii:ss时间格式
func GetFormatYmdHis() string {
    return GetFormatYmdHisByTime(time.Now())
}

// 根据Time获得YY-mm-dd HH:ii:ss时间格式
func GetFormatYmdHisByTime(t time.Time) string {
    return t.Format("2006-01-02 15:04:05")
}

// 根据时间戳获得YY-mm-dd HH:ii:ss时间格式
func GetFormatYmdHisByUnix(u int64) string {
    return GetFormatYmdHisByTime(time.Unix(u, 0))
}

// 获得当前时间戳(单位:秒)
func GetTimer() int64 {
    return GetTimerByTime(time.Now())
}

// 根据Time获得时间戳(单位:秒)
func GetTimerByTime(t time.Time) int64 {
    return t.Unix()
}

// 获得当前时间戳(单位:毫秒)
func GetMicroTimer() int64 {
    return GetMicroTimerByTime(time.Now())
}

// 根据Time获得时间戳(单位:毫秒)
func GetMicroTimerByTime(t time.Time) int64 {
    return int64(t.UnixNano() / 1000000)
}