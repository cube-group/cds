package models

import(
	"sys/core"
	"time"
	"net/http"
	"alex/errors"
    "alex/utils"
    "github.com/asaskevich/govalidator"
)

//fcds.f_jobcmds表模型
type FJobCmd struct {
	ID            uint `gorm:"AUTO_INCREMENT"`
	Tid           uint
	Name          string `gorm:"type:varchar(50)"`
	Value         string `gorm:"type:varchar(100)"`
	CreateTime    time.Time `gorm:"datetime"`
	UpdateTime    time.Time `gorm:"datetime"`
}

//验证结构体
type SetsJobCmdValidate struct {
    Tid      uint 
    Name     string `valid:"required~命令名称必填"`
    Value    string  `valid:"required~ssh账号必填"`
}

//根据结构体批量设置属性
func (t *FJobCmd) setByValidateValue(value *SetsJobCmdValidate) {
    t.Tid = value.Tid
    t.Name = value.Name
    t.Value = value.Value

}
//根据req对象返回Validate结构体
func (t *FJobCmd) NewValidateValue(req *http.Request) *SetsJobCmdValidate {
    return &SetsJobCmdValidate{
        Tid:utils.MustUint(req.FormValue("serverId")),
        Name:req.FormValue("name"),
        Value:req.FormValue("script"),
    }
}

func NewFJobCmd() *FJobCmd {
    return new(FJobCmd)
}

//列表
func (t *FJobCmd) PageList(req *http.Request) (interface{}, errors.IMyError) {
     db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)
    var list []*FJobCmd
    res, err := PageListById(req, &list, db.Order("id DESC"));
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    serverList , _ := NewSetsServer().PageList(req)
    servers := *serverList.(map[string]interface{})["list"].(*[]*SetsServer)
    for _, a := range servers{
        if utils.MustUint(a.ID) ==utils.MustUint(req.FormValue("id")){
            res["server"] = a
        }
    }
    return res, nil    
}

//添加
func (t *FJobCmd) Create(req *http.Request) (interface{}, errors.IMyError) {
    validValue := t.NewValidateValue(req)
    if ok, err := govalidator.ValidateStruct(validValue); !ok || err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    defer core.PoolRollBack(db)

    obj := new(FJobCmd)
    obj.setByValidateValue(validValue)
    obj.CreateTime = time.Now()
    if err := db.Save(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    return map[string]interface{}{"id":obj.ID}, nil
}

//修改
func (t *FJobCmd) Update(req *http.Request) (interface{}, errors.IMyError) {
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
    obj := new(FJobCmd)
    db.Where("id=?", id).First(obj)
    if obj.ID == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "记录不存在！")
    }
    obj.setByValidateValue(validValue)
    obj.UpdateTime = time.Now()
    /*if err := db.Save(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }*/
    if err := db.Model(obj).Update(FJobCmd{Name: obj.Name, Value: obj.Value}).Error ; err !=nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }

    return nil, nil
}

//删除
func (t *FJobCmd) Del(req *http.Request, userInfo *ContextInfo) (interface{}, errors.IMyError) {
    id := utils.MustInt(req.FormValue("id"))
    if id == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "id必传")
    }

    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    obj := new(FJobCmd)
    db.Where("id=?", id).First(obj)
    if obj.ID == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "记录不存在")
    }
    if err := db.Delete(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    return nil, nil
}