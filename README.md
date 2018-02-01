### CDS-Container Deploy System
![](https://github.com/cube-group/cds/blob/master/images/cds-framework.png)
### web framework
martini
### attension
* 未使用session插件(感兴趣的可以看下martini的session插件)
* 每次访问页面都会请求fcds.f_users表进行cookie中的token状态检测
### 目录解释
* conf 配置文件目录
* controllers 路由目录
* core 核心工具类
* data 全局常用info工具类
* log 日志目录
* models 数据代理
* plugins 中间件和插件
* servers 没卵用
* sql 项目sql集合
* views 页面模板
* main.go 程序主入口
### martini全局注入可用实例
* martini.Context martini原始注入(可操作本次请求的依赖注入)
* *http.Request http请求
* http.ResponseWriter http返回
* *core.Pools mysql和redis连接池
* render.Render 返回值和页面渲染
* *log.Logger 日志
* *data.UserInfo 用户信息对象
### 运行sys程序
```
$cd $GOPATH/src/sys && go run main.go
```
### 编译sys程序
```
$cd $GOPATH/src && go build && ./sys
```