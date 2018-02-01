package models

import (
    "sys/core"
    "alex/errors"
    "github.com/asaskevich/govalidator"
    "fmt"
    "time"
    "alex/utils"
    "net/url"
)

const (
    IMAGE_BUILD_ERR = 0 //镜像构建失败
    IMAGE_BUILD_ING = 1 //镜像构建中
    IMAGE_BUILD_COMPLETE = 2 //镜像构建完成
)

//fcds.f_build镜像构建表
type FBuild struct {
    Id         int `gorm:"int(11) AUTO_INCREMENT",json:"id"`         //id
    Username   string `gorm:"varchar(100) not null"`                 //操作人用户名
    UniqueId   string `gorm:"varchar(200) not null",json:"uniqueId"` //编译的唯一id
    BaseImage  string `gorm:"varchar(100) not null"`                 //基础镜像名称(不包括:latest)
    ServiceId  int `gorm:"int(11)",json:"serviceId"`                 //微服务id
    Name       string`gorm:"varchar(100)",json:"name"`               //微服务名称
    Version    string `gorm:"varchar(100)",json:"version"`           //微服务版本号
    V1         int `gorm:"int"`
    V2         int `gorm:"int"`
    V3         int `gorm:"int"`
    Src        string `gorm:"varchar(500)",json:"src"`               //携带打包项目地址
    Path       string `gorm:"varchar(500)",json:"path"`              //容器内项目文件所在目录
    Status     int `gorm:"tinyint(2)",json:"status"`                 //编译状态码
    StatusMsg  string `gorm:"varchar(100)",json:"statusMsg"`         //编译状态
    CreateTime time.Time `gorm:"datetime",json:"createTime"`
    Log        string `gorm:"text",json:"log"`                       //编译日志
}

//build结构体
type FBuildPost struct {
    baseImage string `valid:"required"`     //基础镜像名称(不带版本号,fcds支持latest的基础服务)
    msId      int `valid:"required"`        //微服务id
    msName    string `valid:"required"`     //微服务名称
    msVersion string `valid:"required"`     //微服务版本
    src       string `valid:"url,required"` //跟随打包的项目地址(svn或git)
}

func NewBuildModel() *FBuild {
    return new(FBuild)
}

//获取构建历史
func (this *FBuild)GetHistory(get url.Values) (utils.PageMapData, errors.IMyError) {
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_BUILD_INDEX, "构建历史读取错误", err.Error())
    }
    defer core.PoolRollBack(db)

    var total uint
    err = db.Model(&FBuild{}).Count(&total).Error
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_BUILD_INDEX, "构建历史读取错误", err.Error())
    }

    var list []*FBuild
    pageDetail := utils.Page(utils.StringToUint(get.Get("page")), utils.StringToUint(get.Get("pageSize")), total)
    err = db.Limit(pageDetail.LimitString).Order("id desc").Find(&list).Error
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_BUILD_INDEX, "构建历史读取错误", err.Error())
    }

    return utils.PageMapData{"page":pageDetail, "list":list}, nil
}

//构建微服务
func (this *FBuild)Create(post url.Values, sets *FSets, username string) (*FBuild, errors.IMyError) {
    //参数检测
    postStruct := FBuildPost{
        baseImage:post.Get("baseImage"),
        msId      :utils.StringToInt(post.Get("msId")),
        msName    :post.Get("msName"),
        msVersion :post.Get("msVersion"),
        src       :post.Get("src"),
    }
    ok, err := govalidator.ValidateStruct(postStruct)
    if !ok || err != nil {
        return nil, errors.NewCodeErr(core.ERR_BUILD_CREATE, "参数不合法", err.Error())
    }

    //版本号检测
    v1, v2, v3, err := utils.IsVersion3(postStruct.msVersion)
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_BUILD_CREATE, "版本号不合法", err.Error())
    }

    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_BUILD_CREATE, "无法连接至数据库", err.Error())
    }

    versionFind := new(FBuild)
    sql := "?=name AND ((?<v1) OR (?=v1 AND ?<v2) OR (?=v1 AND ?=v2 AND ?<=v3))"
    if err := db.Where(sql, postStruct.msName, v1, v1, v2, v1, v2, v3).Last(&versionFind).Error; err != nil {
        fmt.Println("暂未发先更高版本可继续构建")
    }
    if (versionFind.Id > 0) {
        return nil, errors.NewCodeErr(core.ERR_BUILD_CREATE, "构建的版本号过低请重新填写,目前最高版本:", versionFind.Version)
    }

    //构建前录库
    model := &FBuild{
        Username:username,
        UniqueId:"",
        BaseImage:postStruct.baseImage,
        ServiceId:postStruct.msId,
        Name:postStruct.msName,
        Version:postStruct.msVersion,
        V1:v1,
        V2:v2,
        V3:v3,
        Src:postStruct.src,
        Path:fmt.Sprintf("/opt/%v", postStruct.msName),
        Status:IMAGE_BUILD_ING,
        StatusMsg:"构建中",
        CreateTime:time.Now(),
    }
    if err := db.Save(model).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_BUILD_CREATE, "构建记录无法创建,无法继续构建镜像", err.Error())
    }

    token, logs, err := buildImage(sets, postStruct)
    if err != nil {
        model.UniqueId = token
        model.Status = IMAGE_BUILD_ERR
        model.StatusMsg = "构建失败"
    } else {
        model.UniqueId = token
        model.Status = IMAGE_BUILD_COMPLETE
        model.StatusMsg = "构建完成"
    }
    model.Log = logs
    if err := db.Save(model).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_BUILD_CREATE, "编译失败,请尝试重新构建", err.Error())
    }
    if model.Status == IMAGE_BUILD_ERR {
        return nil, errors.NewCodeErr(core.ERR_BUILD_CREATE, "编译失败,请尝试重新构建\n", model.Log)
    }

    return model, nil
}

//微服务镜像构建页面返回
func (this *FBuild)CreateIndex(get url.Values, sets *FSets) (map[string]interface{}, errors.IMyError) {
    //获取所有基础镜像
    registry := NewDockerNodeBaseRegistry(sets)
    list, err := registry.GetImages()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_BUILD_INDEX, err.Error(), "无法载入基础镜像列表")
    }
    //获取所有微服务名称
    msList, e := GetAllMs()
    if e != nil {
        return nil, errors.NewCodeErr(core.ERR_MS_LIST, e.String(), "无法载入微服务名称列表")
    }
    return map[string]interface{}{"images":list, "ms":msList}, nil
}

//构建镜像函数
func buildImage(sets *FSets, s FBuildPost) (string, string, error) {
    builder := NewDockerNodeBuild(sets)

    valid, err := builder.IsValid()
    if !valid || err != nil {
        return "", builder.Logs(), err
    }

    out, err := builder.DockerBuildAndPush(
        s.msName,
        s.msVersion,
        fmt.Sprintf("%v:latest", s.baseImage),
        s.src,
    )
    if err != nil {
        return "", builder.Logs(), err
    }
    return out, builder.Logs(), nil
}