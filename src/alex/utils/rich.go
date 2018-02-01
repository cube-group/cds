package utils

import (
    "regexp"
    "errors"
    "strings"
)

func IsMatched(pattern string, value interface{}) bool {
    matched, err := regexp.Match(pattern, value.([]byte))
    if err != nil {
        return false
    }
    return matched
}

// 是否为纯数字
func IsNumber(value interface{}) bool {
    pattern := `^[0-9]*$`
    return IsMatched(pattern, value)
}

// 判断是否为汉字编码
func IsLetterAndNumber(value interface{}) bool {
    pattern := `^[A-Za-z0-9]+$`
    return IsMatched(pattern, value)
}

// 判断是否为汉字编码
func IsChineseEncode(value interface{}) bool {
    pattern := `^[\u4e00-\u9fa5]{0,}$`
    return IsMatched(pattern, value)
}

// 是否为邮箱地址
func IsEmail(value interface{}) bool {
    pattern := `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
    return IsMatched(pattern, value)
}

// 判断字符串是否为url
func IsURL(value interface{}) bool {
    pattern := `^http://(["w-]+\.)+[\w-]+(/[\w-./?%&=]*)?$`
    return IsMatched(pattern, value)
}

// 判断是否为域名
func IsDomain(value interface{}) bool {
    pattern := `[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(/.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+/.?`
    return IsMatched(pattern, value)
}

// 判断是否为手机号
func IsPhone(value interface{}) bool {
    pattern := `^(13[0-9]|14[5|7]|15[0|1|2|3|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\d{8}$`
    return IsMatched(pattern, value)
}

// 判断是否为电话号码
func IsTel(value interface{}) bool {
    pattern := `^($$\d{3,4}-)|\d{3.4}-)?\d{7,8}$`
    return IsMatched(pattern, value)
}

// 判断是否为身份证号码
// 15位、18位数字
// 数字、字母x结尾
func IsID(value interface{}) bool {
    pattern1 := `^\d{15}|\d{18}$`
    pattern2 := `^\d{8,18}|[0-9x]{8,18}|[0-9X]{8,18}?$`
    return IsMatched(pattern1, value) || IsMatched(pattern2, value)
}

// 判断是否为合法账号
// 字母开头，允许5-16字节，允许字母数字下划线
func IsValidUsername(value interface{}) bool {
    pattern := `^[a-zA-Z][a-zA-Z0-9_]{4,15}$`
    return IsMatched(pattern, value)
}

// 判断是否为合法密码
// 以字母开头，长度在6~18之间，只能包含字母、数字和下划线)
func IsValidPassword(value interface{}) bool {
    pattern := `^[a-zA-Z]\w{5,17}$`
    return IsMatched(pattern, value)
}

// 判断是否为合法强效密码
// 必须包含大小写字母和数字的组合，不能使用特殊字符，长度在8-24之间
func IsValidForcePassword(value interface{}) bool {
    pattern := `^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,24}$`
    return IsMatched(pattern, value)
}

// 判断是否为日期格式2017-09-10
func IsDate(value interface{}) bool {
    pattern := `^\d{4}-\d{1,2}-\d{1,2}`
    return IsMatched(pattern, value)
}

//版本号检测如:1.0.0
//返回版本号
func IsVersion3(value string) (int, int, int, error) {
    pattern := `^\d+(\.\d+){2}`
    matched, err := regexp.MatchString(pattern, value)
    if !matched || err != nil {
        return 0, 0, 0, errors.New("version3 format error")
    }
    arr := strings.Split(value, ".")
    return StringToInt(arr[0]), StringToInt(arr[1]), StringToInt(arr[2]), nil
}