package models

import (
    "time"
    "net/url"
    "alex/utils"
    "sys/core"
    "alex/errors"
)

func NewMsModel() *FMs {
    return new(FMs)
}

//微服务表fcds.f_ms表结构
type FMs struct {
    Id         int `gorm:"int AUTO_INCREMENT"`  //微服务部署主键id
    Username   string     `gorm:"varchar(100)"` //发布人用户名称
    MsId       int `gorm:"int"`                 //微服务名称id
    Name       string `gorm:"varchar(100)"`     //微服务名称
    Version    string `gorm:"varchar(100)"`     //微服务版本
    Status     int `gorm:"int"`                 //部署状态
    StatusMsg  string `gorm:"varchar(100)"`     //部署状态文字
    UseTime    int `gorm:"int"`                 //部署耗时
    DeployTye  int `gorm:"int"`                 //部署类型
    FinishTime time.Time `gorm:"datetime"`      //完成时间
    CreateTime time.Time `gorm:"datetime"`      //开始时间
}

//按照微服务名称查看
type DeployDetail struct {
    Id           int    //微服务名称id
    Name         string //微服务名称
    OfficalExist bool   //是否包含正式版本
    Offical      FMs    //正式版本内容
    GrayExist    bool   //是否包含灰度版本
    Gray         []FMs  //灰度版本内容
}

//获取微服务列表
func (this *FMs)PageList(get url.Values) (map[string]interface{}, errors.IMyError) {
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_DEPLOY_INDEX, "无法连接数据库", err.Error())
    }
    defer core.PoolRollBack(db)

    var total uint
    if err = db.Model(FMs{}).Count(&total).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_BUILD_INDEX, "部署列表读取失败", err.Error())
    }

    //读取微服务列表
    data, err := NewSetsMsModel().PageList(get)
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_BUILD_INDEX, "部署列表读取失败", err.Error())
    }

    //针对当前的微服务列表匹配正式版本和上线版本
    var msList []*FSetsMs = data["list"].([]*FSetsMs)
    var list []DeployDetail
    for _, ms := range msList {
        var details []FMs
        db.Where("ms_id=? AND status>?", ms.Id, 0).Find(&details)
        var dd DeployDetail
        if len(details) == 0 {
            dd = DeployDetail{Id:ms.Id, Name:ms.Name}
        } else {
            dd = DeployDetail{Id:ms.Id, Name:ms.Name, Gray:[]FMs{}}
            for _, detail := range details {
                if detail.DeployTye == 1 {
                    //正式
                    dd.OfficalExist = true
                    dd.Offical = detail
                } else {
                    //灰度
                    dd.GrayExist = true
                    dd.Gray = append(dd.Gray, detail)
                }
            }
        }
        list = append(list, dd)
    }
    return utils.PageMapData{"page":data["page"], "list":list}, nil
}

//获取微服务部署详情
func (this *FMs)GetDetailList(get url.Values) (utils.PageMapData, errors.IMyError) {
    msId := get.Get("id");
    if msId == "" {
        return nil, errors.NewCodeErr(core.ERR_DEPLOY_DETAIL_LIST, "微服务id参数缺失")
    }

    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_DEPLOY_DETAIL_LIST, "无法连接数据库", err.Error())
    }
    defer core.PoolRollBack(db)

    var count uint
    if err := db.Model(FMs{}).Where("ms_id=?", msId).Count(&count).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_DEPLOY_DETAIL_LIST, "无法获取分页", err.Error())
    }

    var list[]*FMs
    pageDetail := utils.Page(utils.StringToUint(get.Get("page")), utils.StringToUint(get.Get("pageSize")), count)
    if err := db.Where("ms_id=?", msId).Order("id DESC").Find(&list).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_DEPLOY_DETAIL_LIST, "没有找到任何数据", err.Error())
    }

    return utils.PageMapData{"id":msId, "page":pageDetail, "list":list}, nil
}

//创建微服务部署页面
func (this *FMs)Create(get url.Values) (map[string]interface{}, errors.IMyError) {
    msId := get.Get("id");
    if msId == "" {
        return nil, errors.NewCodeErr(core.ERR_DEPLOY_CREATE, "微服务id参数缺失")
    }

    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_DEPLOY_CREATE, "无法连接数据库", err.Error())
    }
    defer core.PoolRollBack(db)

    var list[]*FMs
    if err := db.Where("status=1 AND service_id=?", msId).Order("id DESC").Limit("0,5").Find(&list).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_DEPLOY_CREATE, "没有可部署版本", err.Error())
    }
    return map[string]interface{}{"Name":list[0].Name, "Versions":list}, nil
}

//正式进行微服务部署
func (this *FMs)Do(post url.Values) {

}