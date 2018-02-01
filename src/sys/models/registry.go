package models

import (
    "github.com/imroc/req"
    "fmt"
    "alex/utils"
    "errors"
)

//仓库Rest处理接口
type IRegistry interface {
    //获取所有镜像列表
    GetImages() ([]interface{}, error)
    //获取某个镜像的所有版本号
    GetImageVersions(name string) ([]string, error)
    //是否包含某个版本的镜像
    HasImage(name, version string) (bool, error)
    //获取所有镜像列表的所有版本号
    GetAllImages() (map[string][]interface{}, error)
    //删除镜像
    //此方法暂时不支持
    //请使用INodeDockerRegistry.DeleteImage进行镜像删除
    DeleteImage(name string, version string) (bool, error)
}

type IRegistryMicroService interface {
    IRegistry

    //确保镜像版本个数
    MakeSureNumbers(name string) (uint, uint, error)
}

//镜像仓库类
type Registry struct {
    Username string
    Password string
    Domain   string
}

func NewRegistry(username, password, domain string) *Registry {
    return &Registry{
        Username:username,
        Password:password,
        Domain:domain,
    }
}

func (this *Registry)request(router, method string, v ...interface{}) (*req.Resp, error) {
    token := fmt.Sprintf("%v:%v", this.Username, this.Password)
    header := req.Header{
        "Accept":        "application/json",
        "Authorization": fmt.Sprintf("Basic %v", utils.Base64Encode(token)),
    }
    url := fmt.Sprintf("https://%v/%v", this.Domain, router)
    var res *req.Resp
    var err error
    if method == "GET" {
        res, err = req.Get(url, header, v)
    } else if method == "POST" {
        res, err = req.Get(url, header, v)
    } else if method == "DELETE" {
        res, err = req.Delete(url, header, v)
    } else if method == "PUT" {
        res, err = req.Put(url, header, v)
    }
    return res, err
}

//获取所有镜像列表
func (this *Registry)GetImages() ([]interface{}, error) {
    var list []interface{}
    res, err := this.request("v2/_catalog", "GET")
    if err != nil {
        return list, err
    }

    var result map[string]interface{}
    err = res.ToJSON(&result)
    if err != nil {
        return list, err
    }
    if (result["repositories"] == nil) {
        return nil, errors.New("数据为空")
    } else {
        return result["repositories"].([]interface{}), nil
    }
}

//获取某个镜像的所有版本号
func (this *Registry)GetImageVersions(name interface{}) ([]interface{}, error) {
    defer func() {
        if e := recover(); e != nil {
            fmt.Println("Registry.GetImageVersions error", e)
        }
    }()

    var list []interface{}
    res, err := this.request(fmt.Sprintf("v2/%v/tags/list", name), "GET")
    if err != nil {
        return list, err
    }

    var result map[string]interface{}
    err = res.ToJSON(&result)
    if err != nil {
        return list, err
    }

    if _, ok := result["tags"]; !ok {
        return list, errors.New("tags is null")
    }

    return result["tags"].([]interface{}), nil
}


//是否包含某个版本号的镜像
func (this *Registry)HasImage(name, version string) (bool, error) {
    versions, err := this.GetImageVersions(name)
    if err != nil {
        return false, err
    }
    if version == "" {
        return true, nil
    }
    for _, v := range versions {
        if v == version {
            return true, nil
        }
    }
    return false, errors.New("image not exists")
}

//获取所有镜像列表的所有版本号
func (this *Registry)GetAllImages() (map[string][]interface{}, error) {
    images, err := this.GetImages()
    if err != nil {
        return map[string][]interface{}{}, err
    }
    var allImages map[string][]interface{} = map[string][]interface{}{}
    for _, image := range images {
        versions, err := this.GetImageVersions(image)
        if err == nil {
            allImages[image.(string)] = versions
        }
    }
    return allImages, nil
}

//删除镜像
//将镜像版本从registry标记删除(逻辑删除)
//物理删除需要操作docker exec中registry的garbage-collect进行操作
//此操作会导致该manifests完全作废
func (this *Registry)DeleteImage(name, version string) (bool, error) {
    token := fmt.Sprintf("%v:%v", this.Username, this.Password)
    header := req.Header{
        "Accept":        "application/vnd.docker.distribution.manifest.v2+json",
        "Authorization": fmt.Sprintf("Basic %v", utils.Base64Encode(token)),
    }

    //获取digest
    url := fmt.Sprintf("https://%v/v2/%v/manifests/%v", this.Domain, name, version)
    res, err := req.Head(url, header)
    if err != nil {
        return false, err
    }

    digest := res.Response().Header.Get("Docker-Content-Digest")
    if digest == "" {
        return false, errors.New("digest not found in header")
    }

    //执行逻辑删除
    url = fmt.Sprintf("https://%v/v2/%v/manifests/%v", this.Domain, name, digest)
    _, err = req.Delete(url, header)
    if err != nil {
        return false, err
    }
    return true, nil
}