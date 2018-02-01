package models

import(
	"sys/core"
	"time"
	"net/http"
    "net/url"
	"alex/errors"
	"alex/utils"
    "github.com/asaskevich/govalidator"
)

//fcds.f_job_task表模型
type FJobTask struct {
	ID            uint `gorm:"AUTO_INCREMENT"`
	Tid           uint 
	Time          string `gorm:"type:varchar(50)"`
	Value         string `gorm:"type:varchar(100)"`
	CreateTime    time.Time `gorm:"datetime"`
	UpdateTime    time.Time `gorm:"datetime"`
}

//验证结构体
type SetsJobTaskValidate struct {
    Tid       string 
    Time      string `valid:"required~任务时间请填写"`
    Value     string  `valid:"required~任务脚本请填写"`
}

//根据结构体批量设置属性
func (t *FJobTask) setByValidateValue(value *SetsJobTaskValidate) {
    t.Tid = utils.MustUint(value.Tid)
    t.Time = value.Time
    t.Value = value.Value
}
//根据req对象返回Validate结构体
func (t *FJobTask) NewValidateValue(req *http.Request) *SetsJobTaskValidate {
    return &SetsJobTaskValidate{
        Tid:req.FormValue("serverId"),
        Time:req.FormValue("time"),
        Value:req.FormValue("script"),
    }
}

func NewFJobTask() *FJobTask {
    return new(FJobTask)
}

//列表
func (t *FJobTask) PageList(req *http.Request) (interface{}, errors.IMyError) {
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)
    var list []*FJobTask
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
func (t *FJobTask) Create(req *http.Request) (interface{}, errors.IMyError) {
    validValue := t.NewValidateValue(req)
    if ok, err := govalidator.ValidateStruct(validValue); !ok || err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    defer core.PoolRollBack(db)

    obj := new(FJobTask)
    obj.setByValidateValue(validValue)
    obj.CreateTime = time.Now()
    if err := db.Save(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    return map[string]interface{}{"id":obj.ID}, nil
}

//修改
func (t *FJobTask) Update(req *http.Request) (interface{}, errors.IMyError) {
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
    obj := new(FJobTask)
    db.Where("id=?", id).First(obj)
    if obj.ID == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "记录不存在！")
    }
    obj.setByValidateValue(validValue)
    obj.UpdateTime = time.Now()
   /* if err := db.Save(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }*/
    if err := db.Model(obj).Update(FJobTask{Time: obj.Time, Value: obj.Value}).Error ; err !=nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }

    return nil, nil
}

//删除
func (t *FJobTask) Del(req *http.Request, userInfo *ContextInfo) (interface{}, errors.IMyError) {
    id := utils.MustInt(req.FormValue("id"))
    if id == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "id必传")
    }

    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_POOL_MYSQL, err)
    }
    defer core.PoolRollBack(db)

    obj := new(FJobTask)
    db.Where("id=?", id).First(obj)
    if obj.ID == 0 {
        return nil, errors.NewCodeErr(core.ERR_OTHER, "记录不存在")
    }
    if err := db.Delete(obj).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    return nil, nil
}

//重启
func (t *FJobTask) Restart(req *http.Request) (interface{}, errors.IMyError){
    data := make(url.Values)
    jobtaskList , _ := NewFJobTask().PageList(req)
    jobtask := *jobtaskList.(map[string]interface{})["list"].(*[]*FJobTask)
    for _, a := range jobtask{
        data["jobtime"] = append(data["jobtime"],a.Time)
        data["jobvalue"] = append(data["jobvalue"],a.Time)
    }
    res, err := http.PostForm("http://192.168.92.128:9000/addjob",data)
    if err != nil {
        return res, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    return nil, nil
}

//停止
func (t *FJobTask) Stop(req *http.Request) (interface{}, errors.IMyError){
    res, err := http.Get("http://192.168.92.128:9000/stop")
    if err != nil {
        return res, errors.NewCodeErr(core.ERR_OTHER, err)
    }
    return nil, nil
}