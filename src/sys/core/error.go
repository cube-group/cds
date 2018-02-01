package core

//登录错误类型枚举
const (
    ERR_LOGIN_PARAMS = 10000
    ERR_LOGIN_NOTFOUND = 1001
    ERR_LOGIN_TOTP = 1002
    ERR_LOGIN_TIMEOUT = 1003
    ERR_LOGIN_SERVICE = 1008
    ERR_LOGIN_OTHER = 1009

    ERR_BUILD_CREATE = 3000 //构建微服务镜像错误
    ERR_BUILD_INDEX = 3001  //构建微服务页面错误

    ERR_DEPLOY_INDEX = 4000 //微服务部署首页列表
    ERR_DEPLOY_DETAIL_LIST = 4100 //微服务部署详情列表
    ERR_DEPLOY_CREATE = 4200 //微服务创建部署页面错误

    ERR_MS_LIST = 6100 //微服务名称列表错误

    ERR_CONFIG_SET_SERVICE = 6510 //核心配置设置失败
    ERR_CONFIG_SET_OTHER = 6519   //核心配置设置其它错误

    ERR_USER_LIST = 6300 //获取用户列表失败
    ERR_POOL_MYSQL = 9001 //mysql连接池错误

    ERR_USER_DETAIL = 6310 //用户详情
    ERR_USER_CREATE = 6220 //添加用户
    ERR_USER_EDIT = 6330 //更新用户
    ERR_USER_DEL = 6340 //删除用户

    ERR_NODE_LIST = 6200

    ERR_OTHER = 9999 //暂定一个通用的开发
)