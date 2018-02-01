package utils

import (
    "encoding/xml"
    "encoding/json"
    "errors"
    "reflect"
    "strconv"
)

// xml转换
// 将struct/map/slice转换为xml字符串
// @param target interface{} 只能为struct实例
// @return string
func XmlEncode(target interface{}) string {
    result, err := xml.Marshal(target)
    if err != nil {
        return ""
    }
    return string(result)
}

// xml转换
// 将xml字符串转换为struct
// @param value string xml字符串
// @param target interface{} 只能为struct
// @return interface{}
func XmlDecode(value string, target interface{}) (bool, error) {
    err := xml.Unmarshal([]byte(value), target)
    if err != nil {
        return false, err
    }
    return true, nil
}

// json转换
// 将struct/map/slice转换为json字符串
// @param target interface{} 为struct/map/slice实例
// @return string
func JsonEncode(target interface{}) (string, error) {
    result, err := json.Marshal(target)
    if err != nil {
        return "", err
    }
    return string(result), nil
}

// json转换
// 将json字符串转换为struct/map/slice
// @param value string json字符串
// @param target struct/map/slice
// @return bool
func JsonDecode(value string, target interface{}) (bool, error) {
    err := json.Unmarshal([]byte(value), target)
    if err != nil {
        return false, err
    }
    return true, nil
}

//从struct指针实例中取出其所有属性和值
func GetStructMapData(s interface{}) (map[string]string, error) {
    if s == nil {
        return nil, errors.New("struct is nil")
    }
    m := map[string]string{}
    v := reflect.ValueOf(s)
    t := reflect.TypeOf(s)
    total := v.NumField()
    for k := 0; k < total; k++ {
        if v.Field(k).Kind() != reflect.String {
            continue
        }
        m[t.Field(k).Name] = v.Field(k).String()
    }
    return m, nil
}

//遍历struct并进行赋值
//目前仅支持string
func SetStructData(s interface{}, m map[string]string) (error) {
    if s == nil {
        return errors.New("struct is nil")
    }
    v := reflect.ValueOf(s).Elem()
    t := reflect.TypeOf(s).Elem()

    for key, value := range m {
        _, owned := t.FieldByName(key)
        if !owned {
            continue
        }
        field := v.FieldByName(key)
        if field.Kind() != reflect.String {
            continue
        }
        field.SetString(value)
    }
    return nil
}

//字符串转为Int
func StringToInt(s string) int {
    in, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        return 0
    }
    return int(in)
}

//字符串转为Uint
func StringToUint(s string) uint {
    in, err := strconv.ParseUint(s, 10, 64)
    if err != nil {
        return 0
    }
    return uint(in)
}