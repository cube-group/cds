package utils

import (
    "net/http"
    "net/url"
    "mime/multipart"
    "math"
    "fmt"
)

//分页专用结构
type PageMapData map[string]interface{}

//分页工具类
type PageDetail struct {
    Index       uint   `json:"index"`       //当前页码
    Size        uint   `json:"size"`        //分页大小
    Total       uint   `json:"total"`       //总条数
    Count       uint   `json:"count"`       //总页数
    Next        uint   `json:"next"`        //下一页
    Pre         uint   `json:"pre"`         //上一页
    Limit       uint   `json:"limit"`       //limit大小
    Offset      uint   `json:"offset"`      //offset大小
    Enabled     bool   `json:"enabled"`     //是否可用分页
    LimitString string `json:"limitString"` //sql limit 0,1
}

//New出分页工具类
func Page(page, pageSize, totalCount uint) PageDetail {
    if page == 0 {
        page = 1
    }
    if pageSize == 0 {
        pageSize = 20
    }
    pageCount := getCount(totalCount, pageSize)
    limit := getLimit(pageSize)
    offset := getOffset(page, pageSize)
    return PageDetail{
        Index:page,
        Size:pageSize,
        Total:totalCount,
        Count:pageCount,
        Enabled:pageCount == 1,
        Next:getNext(page, pageCount),
        Pre:getPre(page),
        Limit:limit,
        Offset:offset,
        LimitString:fmt.Sprintf("%v,%v", limit, offset),
    }
}

//获取offset
func getOffset(page, pageSize uint) uint {
    return (page - 1 ) * pageSize
}

//获取limit
func getLimit(pageSize uint) uint {
    return pageSize
}

//获取上一页
func getPre(page uint) uint {
    return page - 1
}

//获取下一页
func getNext(page, pageCount uint) uint {
    if pageCount > 1 && page >= 0 && page < pageCount {
        return page + 1
    }
    return 0
}

//获取总页数
func getCount(totalCount, pageSize uint) uint {
    return uint(math.Ceil(float64(totalCount) / float64(pageSize)))
}


//req默认解析大小
var reqGetFormDataMaxMemory int64 = 1024 << 10

// Files maps a string key to a list of values.
// It is typically used for query parameters and form values.
// Unlike in the http.Header map, the keys in a Values map
// are case-sensitive.
type Files map[string][]*multipart.FileHeader

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v Files) Get(key string) *multipart.FileHeader {
    if v == nil {
        return nil
    }
    vs := v[key]
    if len(vs) == 0 {
        return nil
    }
    return vs[0]
}

// Set sets the key to value. It replaces any existing
// values.
func (v Files) Set(key string, value *multipart.FileHeader) {
    v[key] = []*multipart.FileHeader{value}
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (v Files) Add(key string, value *multipart.FileHeader) {
    v[key] = append(v[key], value)
}

// Del deletes the values associated with key.
func (v Files) Del(key string) {
    delete(v, key)
}

//解析encrypt=form-data
//从http.Request中取出所有formData格式的post数据
func ReqGetFormData(req *http.Request) (url.Values, Files, error) {
    err := req.ParseMultipartForm(reqGetFormDataMaxMemory)
    if err != nil {
        return nil, nil, err
    }

    return url.Values(req.MultipartForm.Value), Files(req.MultipartForm.File), nil
}

//解析Content-Type=application/x-www-form-urlencoded类型的post数据
func ReqGetEncoding(req *http.Request) url.Values {
    err := req.ParseForm()
    if err != nil {
        return make(url.Values)
    }
    return req.PostForm
}

//解析url query数据
//从http.Request中获取所有get参数
func ReqGetQuery(req *http.Request) url.Values {
    values, err := url.ParseQuery(req.URL.RawQuery)
    if err != nil {
        return make(url.Values)
    }
    return values
}

//解析url query数据
//从http.Request中获取所有get参数
//并且同时解析page和pageSize
func ReqGetPageQuery(req *http.Request) (url.Values, uint, uint) {
    values := ReqGetQuery(req)
    return values, StringToUint(values.Get("page")), StringToUint(values.Get("pageSize"))
}

//cookie写操作
func CookieSet(w http.ResponseWriter, key, value string, maxAge int) {
    http.SetCookie(
        w,
        &http.Cookie{Name:key, Value:value, Path:"/", MaxAge:maxAge},
    )
}