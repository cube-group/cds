package models

import (
    "time"
    "sys/core"
    "net/http"
    "alex/errors"
    "alex/utils"
    "github.com/asaskevich/govalidator"
    "net/url"
)

//机群
type FSetsMs struct {
    Id         int `gorm:"int(11) AUTO_INCREMENT"`
    Name       string `gorm:"varchar(100)"`
    CreateTime time.Time `gorm:"datetime"`
    UpdateTime time.Time `gorm:"datetime"`
}

//字段对应中文
func (t FSetsMs) FieldCn(field string) string {
    labels := map[string]string{
        "name" : "微服务名称",
    }
    if label, ok := labels[field]; ok {
        return label
    }
    return field
}


//自定义表名
func (t *FSetsMs) TableName() string {
    return "f_sets_ms"
}


//验证结构体
type SetsMicroServiceValidate struct {
    Name string `valid:"required~微服务名称必传"`
}

//根据结构体批量设置属性
func (t *FSetsMs) setByValidateValue(value *SetsMicroServiceValidate) {
    t.Name = value.Name
}

//根据req对象返回Validate结构体
func (t *FSetsMs) NewValidateValue(req *http.Request) *SetsMicroServiceValidate {
    return &SetsMicroServiceValidate{
        Name:req.FormValue("name"),
    }
}

func NewSetsMsModel() *FSetsMs {
    return new(FSetsMs)
}


//分页列表
func (t *FSetsMs) PageList(get url.Values) (utils.PageMapData, errors.IMyError) {
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    var total uint
    if err = db.Model(&FSetsMs{}).Count(&total).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "微服务名称列表错误", err.Error())
    }

    var list []*FSetsMs
    pageDetail := utils.Page(utils.StringToUint(get.Get("page")), utils.StringToUint(get.Get("pageSize")), total)
    if err = db.Limit(pageDetail.LimitString).Order("id desc").Find(&list).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "微服务名称列表错误", err.Error())
    }

    return utils.PageMapData{"page":pageDetail, "list":list}, nil
}

//添加
func (t *FSetsMs) Create(req *http.Request) (interface{}, errors.IMyError) {
    validValue := t.NewValidateValue(req)
    if ok, err := govalidator.ValidateStruct(validValue); !ok || err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    defer core.PoolRollBack(db)

    obj := new(FSetsMs)
    obj.setByValidateValue(validValue)
    obj.CreateTime = time.Now()
    if err := db.Save(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    return map[string]interface{}{"id":obj.Id}, nil
}

//删除
func (t *FSetsMs) Del(req *http.Request, userInfo *ContextInfo) (interface{}, errors.IMyError) {
    id := utils.MustInt(req.FormValue("id"))
    if id == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "id 必传")
    }
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    obj := new(FSetsMs)
    db.Where("id=?", id).First(obj)
    if obj.Id == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "微服务不存在")
    }

    if err := validateTotp(req.FormValue("totp"), userInfo); err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }

    if err := db.Delete(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    return nil, nil
}

//获取所有微服务
func GetAllMs() ([]*FSetsMs, errors.IMyError) {
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    var list []*FSetsMs
    if err := db.Find(&list).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_MS_LIST, err.Error(), "微服务名称列表获取失败")
    }
    return list, nil
}
