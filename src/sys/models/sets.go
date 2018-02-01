package models

import (
    "alex/utils"
    "net/url"
    "sys/core"
    "alex/errors"
    "github.com/asaskevich/govalidator"
    "fmt"
    "time"
)

//fcds.f_sets_configs
type FSetsConfigs struct {
    Id         string `gorm:"int(11) not null AUTO_INCREMENT"`
    Key        string `gorm:"varchar(100) not null"`
    Value      string `gorm:"varchar(500) not null"`
    CreateTime time.Time `gorm:"datetime"`
    UpdateTime time.Time `gorm:"datetime"`
    Comment    string `gorm:"varchar(500)"`
}

//核心数据体
type FSets struct {
    //机器ssh账号密码(入口机、代理机、编译机、仓库机)
    SshUsername            string `valid:"required"`
    //机器ssh账号密码
    SshPassword            string `valid:"required"`
    //入口机host
    FacadeHost             string `valid:"required"`
    //入口机ssh地址(如:127.0.0.1:22)
    FacadeAddress          string `valid:"required"`
    //镜像编译机ssh地址
    BuilderAddress         string `valid:"required"`
    //镜像仓库统一账号密码
    RegistryUsername       string `valid:"required"`
    RegistryPassword       string `valid:"required"`
    //基础镜像仓库
    RegistryBaseAddress    string `valid:"required"`
    RegistryBaseDomain     string `valid:"host,required"`
    RegistryBaseHost       string `valid:"required"`
    RegistryBaseRunName    string `valid:"required"`
    //微服务镜像仓库
    RegistryServiceAddress string  `valid:"required"`
    RegistryServiceDomain  string `valid:"host,required"`
    RegistryServiceHost    string `valid:"required"`
    RegistryServiceRunName string `valid:"required"`
    //系统模式
    AppMode                string `valid:"required"`
    //svn配置
    SvnUsername            string
    SvnPassword            string
    SvnDomain              string
    SvnHost                string
    //git配置
    GitUsername            string
    GitPassword            string
    GitDomain              string
    GitHost                string
    //钉钉配置
    DingDing               string `valid:"url,required"`
    //任务机默认镜像
    NodeTaskDefaultImages  string `valid:"required"`
}

func NewSetsModel() *FSetsConfigs {
    return new(FSetsConfigs)
}


//存储核心配置数据
func (this *FSetsConfigs)SetCoreConfig(post url.Values) (errors.IMyError) {
    postStruct := FSets{
        SshUsername:post.Get("SshUsername"),
        SshPassword:post.Get("SshPassword"),
        FacadeHost:post.Get("FacadeHost"),
        FacadeAddress:post.Get("FacadeAddress"),
        BuilderAddress:post.Get("BuilderAddress"),
        RegistryUsername:post.Get("RegistryUsername"),
        RegistryPassword:post.Get("RegistryPassword"),
        RegistryBaseAddress:post.Get("RegistryBaseAddress"),
        RegistryBaseDomain:post.Get("RegistryBaseDomain"),
        RegistryBaseHost:post.Get("RegistryBaseHost"),
        RegistryBaseRunName:post.Get("RegistryBaseRunName"),
        RegistryServiceAddress:post.Get("RegistryServiceAddress"),
        RegistryServiceDomain:post.Get("RegistryServiceDomain"),
        RegistryServiceHost:post.Get("RegistryServiceHost"),
        RegistryServiceRunName:post.Get("RegistryServiceRunName"),
        AppMode:post.Get("AppMode"),
        SvnUsername:post.Get("SvnUsername"),
        SvnPassword :post.Get("SvnPassword"),
        SvnDomain :post.Get("SvnDomain"),
        SvnHost :post.Get("SvnHost"),
        GitUsername:post.Get("GitUsername"),
        GitPassword :post.Get("GitPassword"),
        GitDomain :post.Get("GitDomain"),
        GitHost:post.Get("GitHost"),
        DingDing:post.Get("DingDing"),
        NodeTaskDefaultImages:post.Get("NodeTaskDefaultImages"),
    }

    ok, err := govalidator.ValidateStruct(postStruct)
    if !ok || err != nil {
        return errors.NewCodeErr(core.ERR_CONFIG_SET_SERVICE, "配置项目数据不合法", err.Error())
    }

    db, err := core.Mysql()
    if err != nil {
        return errors.NewCodeErr(core.ERR_CONFIG_SET_SERVICE, "配置项目数据不合法", err.Error())
    }
    defer core.PoolRollBack(db)

    m, err := utils.GetStructMapData(postStruct)
    if err != nil {
        return errors.NewCodeErr(core.ERR_CONFIG_SET_SERVICE, "获取标准结构数据失败", err.Error())
    }

    for key, value := range m {
        result := &FSetsConfigs{Key:key}
        if err := db.Where("`key`=?", result.Key).First(result).Error; err != nil {
        }
        result.Value = value
        if err = db.Save(result).Error; err != nil {
            fmt.Println("key:", key, "update err", err.Error())
        }
    }

    return nil
}

//获取核心配置数据
func (this *FSetsConfigs)GetCoreConfig() (*FSets, errors.IMyError) {
    db, err := core.Mysql()
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_CONFIG_SET_SERVICE, "获取配置失败", err.Error())
    }
    defer core.PoolRollBack(db)

    var sets []*FSetsConfigs
    if err := db.Find(&sets).Error; err != nil {
        return nil, errors.NewCodeErr(core.ERR_CONFIG_SET_SERVICE, "获取配置失败", err.Error())
    }
    m := map[string]string{}
    for _, item := range sets {
        m[item.Key] = item.Value
    }

    instance := &FSets{}
    err = utils.SetStructData(instance, m)
    if err != nil {
        return nil, errors.NewCodeErr(core.ERR_CONFIG_SET_SERVICE, "赋值失败", err.Error())
    }
    return instance, nil
}