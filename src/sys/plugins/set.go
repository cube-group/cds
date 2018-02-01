package plugins

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"sys/models"
)

//中间件
//获取FCDS所有系统配置
func MiddleWareSet() martini.Handler {
	return func(c martini.Context, r render.Render, info *models.ContextInfo) {
		sets, err := models.NewSetsModel().GetCoreConfig()
		if err == nil {
			info.Sets = sets
			return
		}

		r.HTML(200, "error/error403", "获取系统配置失败")
	}
}
