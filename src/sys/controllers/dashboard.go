package controllers

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"sys/core"
	"alex/io"
)

//文章列表路由相关操作
type Dashboard struct{}

func DashboardController() core.RouterHandler {
	c := new(Dashboard)
	return func(r martini.Router) {
		r.Get("", c.ActionDashboard)
	}
}

func (d *Dashboard)ActionDashboard(req *http.Request, r render.Render) {
	io.OutputHtml(r, "Dashboard", "dashboard/index", nil)
}