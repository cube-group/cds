package models

import(
	"sys/core"
	"time"
	"net/http"
	"alex/errors"
    "alex/utils"
    "github.com/asaskevich/govalidator"
)
//fcds.f_job_src表模型
type FJobSrc struct {
	ID            uint `gorm:"AUTO_INCREMENT"`
	Tid           uint
	Src           string `gorm:"type:varchar(50)"`
	Name          string `gorm:"type:varchar(100)"`
	CreateTime    time.Time `gorm:"datetime"`
	UpdateTime    time.Time `gorm:"datetime"`
}

//验证结构体
type SetsJobSrcValidate struct {
    Tid       uint 
    Src       string `valid:"required~仓库地址必填"`
    Name      string `valid:"required~存在名称必填"`
}

//根据结构体批量设置属性
func (t *FJobSrc) setByValidateValue(value *SetsJobSrcValidate) {
    t.Tid = value.Tid
    t.Src = value.Src
    t.Name = value.Name
}
//根据req对象返回Validate结构体
func (t *FJobSrc) NewValidateValue(req *http.Request) *SetsJobSrcValidate {
    return &SetsJobSrcValidate{
        Tid:utils.MustUint(req.FormValue("serverId")),
        Src:req.FormValue("registry"),
        Name:req.FormValue("name"),
    }
}

func NewFJobSrc() *FJobSrc {
    return new(FJobSrc)
}

//列表
func (t *FJobSrc) PageList(req *http.Request) (interface{}, errors.IMyError) {
   db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)
    var list []*FJobSrc
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
func (t *FJobSrc) Create(req *http.Request) (interface{}, errors.IMyError) {
    validValue := t.NewValidateValue(req)
    if ok, err := govalidator.ValidateStruct(validValue); !ok || err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    defer core.PoolRollBack(db)

    obj := new(FJobSrc)
    obj.setByValidateValue(validValue)
    obj.CreateTime = time.Now()
    if err := db.Save(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    return map[string]interface{}{"id":obj.ID}, nil
}

//修改
func (t *FJobSrc) Update(req *http.Request) (interface{}, errors.IMyError) {
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
    obj := new(FJobSrc)
    db.Where("id=?", id).First(obj)
    if obj.ID == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "记录不存在！")
    }
    obj.setByValidateValue(validValue)
    obj.UpdateTime = time.Now()
    /*if err := db.Save(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }*/
    if err := db.Model(obj).Update(FJobSrc{Name: obj.Name, Src: obj.Src}).Error ; err !=nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }

    return nil, nil
}

//删除
func (t *FJobSrc) Del(req *http.Request, userInfo *ContextInfo) (interface{}, errors.IMyError) {
    id := utils.MustInt(req.FormValue("id"))
    if id == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "id必传")
    }

    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    obj := new(FJobSrc)
    db.Where("id=?", id).First(obj)
    if obj.ID == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "记录不存在")
    }
    if err := db.Delete(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    return nil, nil
}