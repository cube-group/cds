package plugins

import (
	"net/http"
	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
)

func MiddleWareError404() martini.Handler {
	return func(req *http.Request, r render.Render) {
		r.HTML(404, "error/error404", req.URL.Path)
	}
}