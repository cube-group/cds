package core

import (
    "github.com/jinzhu/gorm"
    "io"
    alexCore "alex/core"
    "github.com/go-redis/redis"
    "github.com/go-martini/martini"
    "github.com/spf13/viper"
    "reflect"
    "errors"
    _ "github.com/go-sql-driver/mysql"
)

const (
    USER_TOKEN = "FCDS_SESSION" //cookie name
)

//重定义martini.Router
type  RouterHandler func(martini.Router)

//全局pool
var poolGlobalInstance *pools

//pool初始化
func PoolInit() {
    if poolGlobalInstance == nil {
        poolGlobalInstance = &pools{
            mysql:alexCore.NewPool(
                poolMysqlFactory,
                viper.GetInt("mysql.maxIdle"),
                viper.GetInt("mysql.maxOpen"),
                "mysql",
                viper.GetBool("debug"),
            ),
            redis:nil,
            //redis:alexCore.NewPool(poolRedisFactory, 100, 100, "redis", viper.GetBool("debug")),
        }
    }
}

//全局数据库&缓存连接池
type pools struct {
    mysql *alexCore.Pool
    redis *alexCore.Pool
}

//获取mysql *gorm.DB连接实例
func Mysql() (*gorm.DB, error) {
    c, err := poolGlobalInstance.mysql.Get()
    if err != nil {
        return nil, err
    }
    return c.(*gorm.DB), nil
}

//获取redis *redis.Conn连接实例
func Redis() (*redis.Conn, error) {
    c, err := poolGlobalInstance.redis.Get()
    if err != nil {
        return nil, err
    }
    return c.(*redis.Conn), nil
}

//连接回炉至连接池(也可能直接销毁)
func PoolRollBack(i io.Closer) error {
    closerString := reflect.TypeOf(i).String()
    if closerString == "*gorm.DB" {
        return poolGlobalInstance.mysql.Back(i)
    } else if closerString == "*redis.Conn" {
        return poolGlobalInstance.redis.Back(i)
    } else {
        return errors.New("Closer Type Error")
    }
}

//mysql pool factory
func poolMysqlFactory() (io.Closer, error) {
    conn, err := gorm.Open(viper.GetString("mysql.driver"), viper.GetString("mysql.master"))
    if err != nil {
        return nil, err
    }
    return conn, nil
}

//redis pool factory
func poolRedisFactory() (io.Closer, error) {
    conn := redis.NewClient(&redis.Options{
        Addr:viper.GetString("redis.address"),
        Password:viper.GetString("redis.password"),
        DB:viper.GetInt("redis.db"),
    })
    _, err := conn.Ping().Result()
    if err != nil {
        conn.Close()
        return nil, err
    }
    return conn, nil
}