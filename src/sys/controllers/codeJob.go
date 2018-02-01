package controllers

import (
	"github.com/go-martini/martini"
	"sys/core"
	"github.com/martini-contrib/render"
	"net/http"
	"alex/io"
	"sys/models"
)

//任务
type CodeJob struct {
	io.Base
}

func CodeJobController() core.RouterHandler {
	c := new(CodeJob)
	return func(r martini.Router) {
		r.Get("", c.Index)
		
		//任务机任务
		r.Get("/detail", c.Detail)
		r.Post("/jobtaskCreate", c.PostJobTaskCreate)
		r.Post("/jobtaskEdit", c.PostJobTaskEdit)
		r.Post("/jobtaskDel", c.PostJobTaskDel)
		r.Post("/jobtaskRestart", c.PostJobTaskRestart)
		r.Post("/jobtaskStop", c.PostJobTaskStop)
		//任务机命令
		r.Get("/command", c.Command)
		r.Post("/jobcmdCreate", c.PostJobCmdCreate)
		r.Post("/jobcmdEdit", c.PostJobCmdEdit)
		r.Post("/jobcmdDel", c.PostJobCmdDel)
		//任务机代码仓库
		r.Get("/registry", c.Registry)
		r.Post("/jobsrcCreate", c.PostJobSrcCreate)
		r.Post("/jobsrcEdit", c.PostJobSrcEdit)
		r.Post("/jobsrcDel", c.PostJobSrcDel)
	}
}
func (t *CodeJob) Index(req *http.Request,r render.Render) {
	res, _ := models.NewSetsServer().PageList(req)
	io.OutputHtml(r, "任务机列表", "codeJob/index", res)
}

func (t *CodeJob) Detail(req *http.Request,r render.Render) {
	res, _ := models.NewFJobTask().PageList(req)
	io.OutputHtml(r, "任务机任务列表", "codeJob/detail", res)
}

func (t *CodeJob) Command(req *http.Request,r render.Render) {
	res, _ := models.NewFJobCmd().PageList(req)
	io.OutputHtml(r, "命令管理", "codeJob/command", res)
}

func (t *CodeJob) Registry(req *http.Request,r render.Render) {
	res, _ := models.NewFJobSrc().PageList(req)
	io.OutputHtml(r, "仓库管理", "codeJob/registry", res)
}

/********************************************* 任务机命令 *********************************************/
//任务机命令删除
func (t *CodeJob) PostJobCmdDel(req *http.Request, r render.Render, userInfo *models.ContextInfo) {
	res, err := models.NewFJobCmd().Del(req, userInfo)
	t.JsonAuto(r, res, err, "删除成功")
}

//任务机命令更新
func (t *CodeJob) PostJobCmdEdit(req *http.Request, r render.Render) {
	res, err := models.NewFJobCmd().Update(req)
	t.JsonAuto(r, res, err, "更新成功")
}

//任务机命令添加
func (t *CodeJob) PostJobCmdCreate(req *http.Request, r render.Render) {
	res, err := models.NewFJobCmd().Create(req)
	t.JsonAuto(r, res, err, "添加成功")
}

/********************************************* 任务机任务 *********************************************/
//任务机任务删除
func (t *CodeJob) PostJobTaskDel(req *http.Request, r render.Render, userInfo *models.ContextInfo) {
	res, err := models.NewFJobTask().Del(req, userInfo)
	t.JsonAuto(r, res, err, "删除成功")
}

//任务机任务更新
func (t *CodeJob) PostJobTaskEdit(req *http.Request, r render.Render) {
	res, err := models.NewFJobTask().Update(req)
	t.JsonAuto(r, res, err, "更新成功")
}

//任务机任务添加
func (t *CodeJob) PostJobTaskCreate(req *http.Request, r render.Render) {
	res, err := models.NewFJobTask().Create(req)
	t.JsonAuto(r, res, err, "添加成功")
}
//任务机任务重启
func (t *CodeJob) PostJobTaskRestart(req *http.Request, r render.Render) {
	res, err := models.NewFJobTask().Restart(req)
	t.JsonAuto(r, res, err, "任务重启成功")
}
//任务机任务停止
func (t *CodeJob) PostJobTaskStop(req *http.Request, r render.Render) {
	res, err := models.NewFJobTask().Stop(req)
	t.JsonAuto(r, res, err, "任务停止成功")
}

/********************************************* 任务机仓库 *********************************************/
//任务机命令删除
func (t *CodeJob) PostJobSrcDel(req *http.Request, r render.Render, userInfo *models.ContextInfo) {
	res, err := models.NewFJobSrc().Del(req, userInfo)
	t.JsonAuto(r, res, err, "删除成功")
}

//任务机命令更新
func (t *CodeJob) PostJobSrcEdit(req *http.Request, r render.Render) {
	res, err := models.NewFJobSrc().Update(req)
	t.JsonAuto(r, res, err, "更新成功")
}

//任务机命令添加
func (t *CodeJob) PostJobSrcCreate(req *http.Request, r render.Render) {
	res, err := models.NewFJobSrc().Create(req)
	t.JsonAuto(r, res, err, "添加成功")
}


