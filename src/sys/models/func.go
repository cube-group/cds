package models

import (
    "alex/utils"
    "net/http"
    "errors"
    "github.com/jinzhu/gorm"
)

//分页统一输出
func PageList(req *http.Request, list interface{}, query *gorm.DB) (map[string]interface{}, error) {

    var totalCount uint
    query.Model(list).Count(&totalCount)

    p := NewPageDetail(req, totalCount)

    if err := query.Limit(p.Limit).Offset(p.Offset).Find(list).Error; err != nil {
        return nil, err
    }
    return map[string]interface{}{"page":p, "list":list}, nil
}


//获取分页对象
func NewPageDetail(req *http.Request, totalCount uint) utils.PageDetail {
    return utils.Page(PageIndex(req), PageSize(req), totalCount)
}


//获取页码
func PageIndex(req *http.Request) uint {
    page := utils.MustUint(req.FormValue("page"))
    if page == 0 {
        page = 1
    }
    return page
}

//获取分页大小
func PageSize(req *http.Request) uint {
    pageSize := utils.MustUint(req.FormValue("pageSize"))
    if pageSize == 0 {
        pageSize = 30
    }
    return pageSize
}

//验证totp
func validateTotp(totp string, userInfo *ContextInfo) error {
	if totp == "" {
		return errors.New("动态key错误")
	}
	equal, err := utils.TotpCode(totp, userInfo.User.TotpSecret)

	if err != nil || !equal {
		return errors.New("动态key错误")
	}
	return nil
}

//根据参数id分页统一输出
func PageListById(req *http.Request, list interface{}, query *gorm.DB) (map[string]interface{}, error) {

    var totalCount uint
    query.Model(list).Count(&totalCount)
    p := NewPageDetail(req, totalCount)

    if err := query.Where("Tid = ?",utils.MustUint(req.FormValue("id"))).Limit(p.Limit).Offset(p.Offset).Find(list).Error; err != nil {
        return nil, err
    }
    return map[string]interface{}{"page":p, "list":list}, nil
}