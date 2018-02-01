package main

import (
    "github.com/go-martini/martini"
    "github.com/spf13/viper"
    "fmt"
    "github.com/martini-contrib/render"
    c "sys/controllers"
    p "sys/plugins"
    "sys/core"
    "runtime"
)

func main() {
    //max cpu used
    fmt.Println("NUM-CPU:", runtime.NumCPU())
    runtime.GOMAXPROCS(1)

    //config init
    viper.SetConfigName("init") // name of config file (without extension)
    viper.AddConfigPath("conf") // path to look for the config file in
    viper.SetConfigType("json")
    err2 := viper.ReadInConfig()
    if err2 != nil {
        panic(fmt.Errorf("Fatal error config file: %s \n", err2))
    }

    //pool init
    core.PoolInit()

    //web server init
    m := martini.Classic()
    p.PluginLog(m)
    m.Use(render.Renderer(render.Options{
        Directory: "views", // Specify what path to load the templates from.
        Layout:    "layout/main", // Specify a layout template.
        // Layouts can call {{ yield }} to render the current template.
        Extensions:      []string{".tmpl", ".html"}, // Specify extensions to load for templates.
        Charset:         "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
        IndentJSON:      true, // Output human readable JSON
        IndentXML:       true, // Output human readable XML
        HTMLContentType: "text/html", // Output XHTML content type instead of default "text/html"
    }))

    //middleWare init
    m.Use(p.MiddleWareError500())
    m.NotFound(p.MiddleWareError404())

    //routes init
    initRoutes(m)

    //martini run
    m.Run()
}

//initialize routes
func initRoutes(m *martini.ClassicMartini) {
    m.Group("/user", c.UserController())

    m.Group("/build", c.BuildController(), p.MiddleWareAuth(), p.MiddleWareSet(), p.UserOperationLog())
    m.Group("/deploy", c.DeployController(), p.MiddleWareAuth(), p.MiddleWareSet(), p.UserOperationLog())
    m.Group("/sets", c.SetsController(), p.MiddleWareAuth(), p.MiddleWareSet(), p.UserOperationLog())
    m.Group("/codeJob", c.CodeJobController(), p.MiddleWareAuth(), p.MiddleWareSet(), p.UserOperationLog())
    m.Group("/dashboard", c.DashboardController(), p.MiddleWareAuth(), p.MiddleWareSet(), p.UserOperationLog())
    m.Get("", func(r render.Render) {
        r.Redirect("/dashboard")
    })
    m.Get("/", func(r render.Render) {
        r.Redirect("/dashboard")
    })
}
