package io

import (
	"github.com/martini-contrib/render"
	"alex/errors"
)

//利用martini框架的Render中间件进行渲染json
func OutputJson(r render.Render, data interface{}, code int, msg string) {
	m := map[string]interface{}{
		"data":data,
		"code":code,
		"msg":msg,
	}
	r.JSON(200, m)
}

//利用martini框架的Render中间件进行渲染json
func OutputHtml(r render.Render, title, name string, data interface{}) {
	r.HTML(200, name, map[string]interface{}{"Title":title, "data":data})
}


//基类控制器,封装公共方法
type Base struct {
}

//输出成功或失败json
func (t *Base) JsonAuto(r render.Render, data interface{}, err errors.IMyError, successMsg string) {
	if err != nil {
		OutputJson(r, data, err.Code(), err.Msg())
	} else {
		OutputJson(r, data, 0, successMsg)
	}
}


