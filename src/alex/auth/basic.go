// basic auth基础验证
package auth

import "alex/utils"

// 默认secret
var BasicAuthSecret string = ""

// 校验sign
func BasicAuthCheckSign(sign string, m map[string]interface{}) bool {
	md5Str := BasicAuthGetSign(m)
	if md5Str == "" {
		return false
	}
	if md5Str != sign {
		return false
	}
	return true
}

// 生成sign
// values需要按照key进行正向排序
// 组装出key1=value1&key2=value2...的字符串,这里定义为queryString
// sign = md5(queryString+"&secret="+basicAuthSecret)
func BasicAuthGetSign(m map[string]interface{}) string {
	if m == nil {
		return ""
	}
	query := utils.MapToQuery(m, true)
	if query == "" {
		return ""
	}
	return utils.MD5(query + "&secret=" + BasicAuthSecret)
}