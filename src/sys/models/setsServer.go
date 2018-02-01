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
type SetsServer struct {
    ID              int `gorm:"int(11) AUTO_INCREMENT"`
    Ip              string `gorm:"varchar(200) not null"`
    Port            int `gorm:"int(11)"`
    Type            int `gorm:"tinyint(11)"`
    Images          string `gorm:"varchar(500)"`
    MicroServiceNum int `gorm:"int(11)"`
    Username        string `gorm:"varchar(100)"`
    Password        string `gorm:"varchar(100)"`
    CreateTime      time.Time `gorm:"datetime"`
    UpdateTime      time.Time `gorm:"datetime"`
}

//字段对应中文
func (t SetsServer) FieldCn(field string) string {
    labels := map[string]string{
        "ip" : "IP地址",
        "username" : "ssh账号",
        "password" : "ssh密码",
        "port" : "ssh端口",
        "type" : "机器类型",
        "images" : "镜像",
    }
    if label, ok := labels[field]; ok {
        return label
    }
    return field
}

//自定义表名
func (t *SetsServer) TableName() string {
    return "f_sets_servers"
}

//验证结构体
type SetsServerValidate struct {
    Ip       string `valid:"required~IP必传,ip~IP地址无效"`
    Port     string `valid:"required~端口号必传,port~端口号必传"`
    Username string  `valid:"required~ssh账号必传"`
    Password string  `valid:"required~ssh密码必传"`
    Type     string  `valid:"required~请选择机器类型,range(0|2)~机器类型非法"`
}

//根据结构体批量设置属性
func (t *SetsServer) setByValidateValue(value *SetsServerValidate) {
    t.Type = utils.MustInt(value.Type)
    t.Port = utils.MustInt(value.Port)
    t.Ip = value.Ip
    t.Username = value.Username
    t.Password = value.Password
}

//根据req对象返回Validate结构体
func (t *SetsServer) NewValidateValue(req *http.Request) *SetsServerValidate {
    return &SetsServerValidate{
        Ip:req.FormValue("ip"),
        Port:req.FormValue("port"),
        Type:req.FormValue("type"),
        Username:req.FormValue("username"),
        Password:req.FormValue("password"),
    }
}

func NewSetsServer() *SetsServer {
    return new(SetsServer)
}


//获取类型的中文
func (t *SetsServer) GetTypeCn() string {
    var cn string
    switch t.Type {
    case 0:
        cn = "无效"
    case 1 :
        cn = "生产"
    case 2 :
        cn = "任务"
    default:
        cn = "未知"
    }
    return cn
}

//列表
func (t *SetsServer) PageList(req *http.Request) (interface{}, errors.IMyError) {
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    var list []*SetsServer
    res, err := PageList(req, &list, db.Order("id DESC"));
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    return res, nil
}

//添加
func (t *SetsServer) Create(req *http.Request) (interface{}, errors.IMyError) {
    validValue := t.NewValidateValue(req)
    if ok, err := govalidator.ValidateStruct(validValue); !ok || err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    defer core.PoolRollBack(db)

    obj := new(SetsServer)
    obj.setByValidateValue(validValue)
    obj.CreateTime = time.Now()
    if err := db.Save(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    return map[string]interface{}{"id":obj.ID}, nil
}

//修改
func (t *SetsServer) Update(req *http.Request) (interface{}, errors.IMyError) {
    id := utils.MustInt(req.FormValue("id"))
    validValue := t.NewValidateValue(req)
    if ok, err := govalidator.ValidateStruct(validValue); !ok || err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)
    obj := new(SetsServer)
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
func (t *SetsServer) Del(req *http.Request, userInfo *ContextInfo) (interface{}, errors.IMyError) {
    id := utils.MustInt(req.FormValue("id"))
    if id == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "id必传")
    }

    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    obj := new(SetsServer)
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
