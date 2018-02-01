package controllers

import (
    "github.com/go-martini/martini"
    "sys/core"
    "github.com/martini-contrib/render"
    "net/http"
    "sys/models"
    "alex/utils"
    "alex/io"
)

//镜像构建
type Build struct {
    io.Base
}

func BuildController() core.RouterHandler {
    c := new(Build)
    return func(r martini.Router) {
        r.Get("", c.BuildList)
        r.Get("/create", c.BuildCreateGet)
        r.Post("/create", c.BuildCreatePost)
    }
}

//构建列表
func (t *Build) BuildList(r render.Render, req *http.Request) {
    data, _ := models.NewBuildModel().GetHistory(utils.ReqGetQuery(req))
    io.OutputHtml(r, "微服务构建历史", "build/index", data)
}

//构建镜像页面
func (t *Build) BuildCreateGet(r render.Render, req *http.Request, info *models.ContextInfo) {
    data, _ := models.NewBuildModel().CreateIndex(utils.ReqGetEncoding(req), info.Sets)
    io.OutputHtml(r, "微服务构建", "build/create", data)
}

//构建镜像动作
func (t *Build) BuildCreatePost(r render.Render, req *http.Request, info *models.ContextInfo) {
    data, err := models.NewBuildModel().Create(utils.ReqGetEncoding(req), info.Sets, info.User.Username)
    if err != nil {
        io.OutputJson(r, data, err.Code(), err.Msg())
    } else {
        io.OutputJson(r, data, 0, "构建成功")
    }
}
