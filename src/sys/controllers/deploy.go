package controllers

import (
    "github.com/go-martini/martini"
    "sys/core"
    "github.com/martini-contrib/render"
    "net/http"
    "alex/io"
    "sys/models"
    "alex/utils"
)

//部署
type Deploy struct {
    io.Base
}

func DeployController() core.RouterHandler {
    c := new(Deploy)
    return func(r martini.Router) {
        r.Get("", c.Index)
        r.Get("/do", c.Do)
        r.Get("/create", c.Create)
        r.Get("/detailList", c.DetailList)
    }
}

//微服务部署首页
func (t *Deploy) Index(req *http.Request, r render.Render) {
    data, _ := models.NewMsModel().PageList(utils.ReqGetEncoding(req))
    io.OutputHtml(r, "部署列表", "deploy/index", data)
}

//创建新部署详情页
func (t *Deploy) Create(req *http.Request, r render.Render) {
    data, _ := models.NewMsModel().Create(utils.ReqGetEncoding(req))
    io.OutputHtml(r, "执行部署", "deploy/create", data)
}

//执行微服务部署
func (this *Deploy)Do(req *http.Request, r render.Render) {

}

//微服务部署详情页面
func (t *Deploy) DetailList(req *http.Request, r render.Render) {
    data, _ := models.NewMsModel().GetDetailList(utils.ReqGetQuery(req))
    io.OutputHtml(r, "部署详情", "deploy/detail-list", data)
}
