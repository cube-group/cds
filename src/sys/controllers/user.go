package controllers

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"sys/models"
	"github.com/spf13/viper"
	"alex/utils"
	"alex/io"
	"sys/core"
	"sys/plugins"
)

//用户管理
type User struct {
	io.Base
}

func UserController() core.RouterHandler {
	c := new(User)
	return func(r martini.Router) {
		r.Get("/login", plugins.MiddleWareUserLogin(), c.ActionLoginGet)
		r.Post("/login", c.ActionLoginPost)
		r.Get("/logout", c.ActionLogout)
		r.Get("/register", c.ActionRegisterPost)
	}
}

//用户登录页
func (user *User) ActionLoginGet(r render.Render) {
	r.HTML(200, "user/login", map[string]interface{}{"Title": "登录", }, render.HTMLOptions{Layout: ""})
}

//用户登录请求
func (user *User) ActionLoginPost(req *http.Request, w http.ResponseWriter, r render.Render) {
	token, err := models.NewUserModel().Login(utils.ReqGetEncoding(req))
	if err != nil {
		io.OutputJson(r, nil, err.Code(), err.Msg())
	} else {
		utils.CookieSet(w, core.USER_TOKEN, token, viper.GetInt("session.maxAge"))
		io.OutputJson(r, nil, 0, "登录成功")
	}
}

//用户退出
func (user *User) ActionLogout(w http.ResponseWriter, r render.Render) {
	utils.CookieSet(w, core.USER_TOKEN, "", -1)
	r.Redirect("/user/login")
}

//用户注册请求
func (user *User) ActionRegisterPost(req *http.Request, r render.Render) {
	r.Text(200, "adsfaf")
}
