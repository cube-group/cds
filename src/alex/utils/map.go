// map相关处理
package utils

import (
    "sort"
    "net/url"
    "strings"
    "fmt"
    "reflect"
    "errors"
)

// url.Values排序返回最新的url.Values
// 由于go lang的map特性遍历map时key会进行随机打印
// 所以我们只能做到返回经过排序的map key数组
func ValuesSort(m url.Values) []string {
    var sortKeyArr []string
    for key, _ := range m {
        sortKeyArr = append(sortKeyArr, key)
    }
    sort.Strings(sortKeyArr)

    return sortKeyArr
}

// url.Values转换为query string
func ValuesToQuery(v url.Values, sortFlag bool) string {
    if v == nil {
        return ""
    }

    var queryArr []string
    if sortFlag {
        keyMap := ValuesSort(v)
        for _, newKey := range keyMap {
            queryArr = append(queryArr, newKey + "=" + v.Get(newKey))
        }
    } else {
        for key, _ := range v {
            queryArr = append(queryArr, key + "=" + v.Get(key))
        }
    }
    return strings.Join(queryArr, "&")
}

// query string或者url地址转换为url.Values
func ValuesFromQuery(value string) url.Values {
    if !IsURL(value) {
        value = "http://www.a.com?" + value
    }
    var urlValues url.Values
    urlInstance, err := url.Parse(value)
    if err != nil {
        return nil
    } else {
        urlValues = urlInstance.Query()
    }

    return urlValues
}

// map排序返回最新的map
// 由于go lang的map特性遍历map时key会进行随机打印
// 所以我们只能做到返回经过排序的map key数组
func MapSort(m map[string]interface{}) []string {
    var sortKeyArr []string
    for key, _ := range m {
        sortKeyArr = append(sortKeyArr, key)
    }
    sort.Strings(sortKeyArr)

    return sortKeyArr
}

// map转换为query string
func MapToQuery(m map[string]interface{}, sortFlag bool) string {
    if m == nil {
        return ""
    }

    var queryArr []string
    if sortFlag {
        keyMap := MapSort(m)
        for _, newKey := range keyMap {
            queryArr = append(queryArr, fmt.Sprintf("%v=%v", newKey, m[newKey]))
        }
    } else {
        for key, item := range m {
            queryArr = append(queryArr, fmt.Sprintf("%v=%v", key, item))
        }
    }

    return strings.Join(queryArr, "&")
}

// query string或者url地址转换为map
func MapFromQuery(value string) map[string]interface{} {
    urlValues := ValuesFromQuery(value)
    if urlValues == nil {
        return nil
    }

    newMap := map[string]interface{}{}
    for key, _ := range urlValues {
        newMap[key] = urlValues.Get(key)
    }

    return newMap
}

//map的值映射到struct的对象当中
//仅支持struct值全为string类型
func MapToStruct(i interface{}, m map[string]string) error {
    t := reflect.TypeOf(i).Kind()
    if t != reflect.Struct {
        return errors.New("i type is not struct")
    }

    v := reflect.ValueOf(i).Elem()
    for key, value := range m {
        field := v.FieldByName(key)
        if field.IsNil() || field.IsValid() {
            continue
        }
        if field.Kind() != reflect.String {
            continue
        }
        field.SetString(value)
    }
    return nil
}
