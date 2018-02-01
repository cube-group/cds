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

//配置
type Sets struct {
	io.Base
}

func SetsController() core.RouterHandler {
	c := new(Sets)
	return func(r martini.Router) {
		//用户管理
		r.Get("/userGet", c.GetUser)
		r.Get("/users", c.GetUsers)
		r.Post("/userCreate", c.PostUserCreate)
		r.Post("/userEdit", c.PostUserEdit)
		r.Post("/userDel", c.PostUserDel)
		r.Post("/totpImage", c.PostTotpImage)

		//机群管理
		r.Get("/nodes", c.GetNodes)
		r.Post("/nodeCreate", c.PostNodeCreate)
		r.Post("/nodeEdit", c.PostNodeEdit)
		r.Post("/nodeDel", c.PostNodeDel)

		//微服务设置
		r.Get("/microService", c.GetMicroServices)
		r.Post("/microCreate", c.PostMicroCreate)
		r.Post("/microDel", c.PostMicroDel)

		//代理设置
		r.Get("/proxy", c.GetProxy)
		r.Post("/proxyCreate", c.PostProxyCreate)
		r.Post("/proxyEdit", c.PostProxyEdit)
		r.Post("/proxyDel", c.PostProxyDel)

		//核心设置
		r.Get("/core", c.CoreGet)
		r.Post("/core", c.CoreSet)
	}
}


/********************************************* 用户设置 *********************************************/
//用户删除
func (t *Sets) PostUserDel(req *http.Request, r render.Render, userInfo *models.ContextInfo) {
	res, err := models.NewUserModel().Del(req, userInfo)
	t.JsonAuto(r, res, err, "删除成功")
}

//用户更新
func (t *Sets) PostUserEdit(req *http.Request, r render.Render) {
	res, err := models.NewUserModel().Update(req)
	t.JsonAuto(r, res, err, "更新成功")
}

//用户添加
func (t *Sets) PostUserCreate(req *http.Request, r render.Render) {
	res, err := models.NewUserModel().Create(req)
	t.JsonAuto(r, res, err, "添加成功")
}

//用户详情
func (t *Sets) GetUser(req *http.Request, r render.Render) {
	res, err := models.NewUserModel().Get(req)
	t.JsonAuto(r, res, err, "获取详情成功")
}

//用户列表
func (t *Sets) GetUsers(req *http.Request, r render.Render) {
	res, _ := models.NewUserModel().PageList(req)
	io.OutputHtml(r, "用户列表", "sets/users", res)
}

//生成用户二维码图片并返回图片文件地址path
func (t *Sets) PostTotpImage(req *http.Request, r render.Render) {
	res, err := models.NewUserModel().TotpImage(req)
	t.JsonAuto(r, res, err, "获取图片")
}

/********************************************* 机群设置 *********************************************/
func (t *Sets) GetNodes(req *http.Request, r render.Render) {
	res, _ := models.NewSetsServer().PageList(req)
	io.OutputHtml(r, "机群管理", "sets/nodes", res)
}

func (t *Sets) PostNodeDel(req *http.Request, r render.Render, userInfo *models.ContextInfo) {
	res, err := models.NewSetsServer().Del(req, userInfo)
	t.JsonAuto(r, res, err, "删除成功")
}

func (t *Sets) PostNodeEdit(req *http.Request, r render.Render) {
	res, err := models.NewSetsServer().Update(req)
	t.JsonAuto(r, res, err, "更新成功")
}

func (t *Sets) PostNodeCreate(req *http.Request, r render.Render) {
	res, err := models.NewSetsServer().Create(req)
	t.JsonAuto(r, res, err, "添加成功")
}


/********************************************* 代理设置 *********************************************/
func (t *Sets) GetProxy(req *http.Request, r render.Render) {
	res, _ := models.NewSetsProxy().PageList(req)
	io.OutputHtml(r, "代理设置", "sets/proxy", res)
}

func (t *Sets) PostProxyDel(req *http.Request, r render.Render, userInfo *models.ContextInfo) {
	res, err := models.NewSetsProxy().Del(req, userInfo)
	t.JsonAuto(r, res, err, "删除成功")
}

func (t *Sets) PostProxyEdit(req *http.Request, r render.Render) {
	res, err := models.NewSetsProxy().Update(req)
	t.JsonAuto(r, res, err, "更新成功")
}

func (t *Sets) PostProxyCreate(req *http.Request, r render.Render) {
	res, err := models.NewSetsProxy().Create(req)
	t.JsonAuto(r, res, err, "添加成功")
}


/********************************************* 微服务设置 *********************************************/
func (t *Sets) GetMicroServices(req *http.Request, r render.Render) {
	res, _ := models.NewSetsMsModel().PageList(utils.ReqGetQuery(req))
	io.OutputHtml(r, "微服务设置", "sets/microService", res)
}

func (t *Sets) PostMicroDel(req *http.Request, r render.Render, userInfo *models.ContextInfo) {
	res, err := models.NewSetsMsModel().Del(req, userInfo)
	t.JsonAuto(r, res, err, "删除成功")
}

func (t *Sets) PostMicroCreate(req *http.Request, r render.Render) {
	res, err := models.NewSetsMsModel().Create(req)
	t.JsonAuto(r, res, err, "添加成功")
}


/********************************************* 核心设置 *********************************************/
func (t *Sets) CoreGet(r render.Render) {
	result, _ := models.NewSetsModel().GetCoreConfig()
	io.OutputHtml(r, "核心设置", "sets/core", result)
}

func (t *Sets) CoreSet(r render.Render, req *http.Request) {
	err := models.NewSetsModel().SetCoreConfig(utils.ReqGetEncoding(req))
	t.JsonAuto(r, nil, err, "修改成功")
}


