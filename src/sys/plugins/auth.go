package plugins

import (
    "github.com/martini-contrib/render"
    "github.com/go-martini/martini"
    "net/http"
    "sys/core"
    "sys/models"
    "alex/io"
)

//中间件
//获取用户登录状态
func MiddleWareAuth() martini.Handler {
    return func(context martini.Context, req *http.Request, w http.ResponseWriter, r render.Render) {
        token, err := models.UserToken(req)
        if err == nil {
            user, err := models.GetUserInfo(token)
            if err == nil {
                context.Map(&models.ContextInfo{User:user})
                return
            }
        }

        //utils.CookieSet(w, core.USER_TOKEN, "", -1)
        if (req.Method == http.MethodGet) {
            r.Redirect("/user/login")
        } else {
            io.OutputJson(r, nil, core.ERR_LOGIN_TIMEOUT, "登录状态超时")
        }
    }
}

//中间件
//user/login专用
func MiddleWareUserLogin() martini.Handler {
    return func(req *http.Request, r render.Render) {
        token, err := models.UserToken(req)
        if err == nil {
            _, err := models.GetUserInfo(token)
            if err == nil {
                r.Redirect("/dashboard")
                return
            }
        }
    }
}

//中间件
//记录用户操作日志
func UserOperationLog() martini.Handler {
    return func(context martini.Context, req *http.Request) {
        token, err := models.UserToken(req)
        if err == nil {
            user, _ := models.GetUserInfo(token)
            log := models.NewFLog()
            log.Uid = user.ID
            log.Username = user.Username
            log.Route = req.RequestURI
            models.Create(log)
        }

    }
}