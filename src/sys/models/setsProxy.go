package models

import (
    "time"
    "sys/core"
    "net/http"
    "alex/errors"
    "alex/utils"
    "github.com/asaskevich/govalidator"
)

//机群
type SetsProxy struct {
    ID         int `gorm:"int(11) AUTO_INCREMENT"`
    Ip         string `gorm:"varchar(200) not null"`
    Type       int `gorm:"tinyint(2)"`
    Port       int `gorm:"int(11)"`
    Status     int `gorm:"tinyint(2)"`
    StatusMsg  string `gorm:"varchar(100)"`
    CreateTime time.Time `gorm:"datetime"`
    UpdateTime time.Time `gorm:"datetime"`
}

//字段对应中文
func (t SetsProxy) FieldCn(field string) string {
    labels := map[string]string{
        "ip" : "IP地址",
        "port" : "端口号",
        "type" : "类型",
    }
    if label, ok := labels[field]; ok {
        return label
    }
    return field
}

//自定义表名
func (t *SetsProxy) TableName() string {
    return "f_sets_proxys"
}

//验证结构体
type SetsProxyValidate struct {
    Ip   string `valid:"required~IP必传,ip~IP地址无效"`
    Port string `valid:"required~端口号必传,port~端口号必传"`
    Type string  `valid:"required~请选择类型,range(0|2)~类型非法"`
}

//根据结构体批量设置属性
func (t *SetsProxy) setByValidateValue(value *SetsProxyValidate) {
    t.Type = utils.MustInt(value.Type)
    t.Port = utils.MustInt(value.Port)
    t.Ip = value.Ip
}

//根据req对象返回Validate结构体
func (t *SetsProxy) NewValidateValue(req *http.Request) *SetsProxyValidate {
    return &SetsProxyValidate{
        Ip:req.FormValue("ip"),
        Port:req.FormValue("port"),
        Type:req.FormValue("type"),
    }
}

func NewSetsProxy() *SetsProxy {
    return new(SetsProxy)
}

//获取类型的中文
func (t *SetsProxy) GetTypeCn() string {
    var cn string
    switch t.Type {
    case 0:
        cn = "Nginx"
    case 1 :
        cn = "SpringCloud"
    case 2 :
        cn = "GoKit"
    default:
        cn = "未知"
    }
    return cn
}


//分页列表
func (t *SetsProxy) PageList(req *http.Request) (interface{}, errors.IMyError) {
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    var list []*SetsProxy
    res, err := PageList(req, &list, db.Order("id DESC"));
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    return res, nil
}

//添加
func (t *SetsProxy) Create(req *http.Request) (interface{}, errors.IMyError) {
    validValue := t.NewValidateValue(req)
    if ok, err := govalidator.ValidateStruct(validValue); !ok || err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    defer core.PoolRollBack(db)

    obj := new(SetsProxy)
    obj.setByValidateValue(validValue)
    obj.CreateTime = time.Now()
    if err := db.Save(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    return map[string]interface{}{"id":obj.ID}, nil
}

//修改
func (t *SetsProxy) Update(req *http.Request) (interface{}, errors.IMyError) {
    id := utils.MustInt(req.FormValue("id"))
    if id == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "id必传")
    }
    validValue := t.NewValidateValue(req)
    if ok, err := govalidator.ValidateStruct(validValue); !ok || err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)
    obj := new(SetsProxy)
    db.Where("id=?", id).First(obj)
    if obj.ID == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "记录不存在！")
    }
    obj.setByValidateValue(validValue)
    obj.UpdateTime = time.Now()
    if err := db.Save(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    return nil, nil
}

//删除
func (t *SetsProxy) Del(req *http.Request, userInfo *ContextInfo) (interface{}, errors.IMyError) {
    id := utils.MustInt(req.FormValue("id"))
    if id == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "id必传")
    }
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    obj := new(SetsProxy)
    db.Where("id=?", id).First(obj)
    if obj.ID == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "记录不存在")
    }
    if err := validateTotp(req.FormValue("totp"), userInfo); err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }

    if err := db.Delete(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    return nil, nil
}


//获取状态的中文
func (t *SetsProxy) GetStatusCn() string {
    var cn string
    switch t.Status {
    case 0:
        cn = "在线"
    case 1 :
        cn = "断开"
    default:
        cn = "未知"
    }
    return cn
}
