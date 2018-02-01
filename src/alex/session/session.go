package session

import (
	"github.com/go-redis/redis"
	"time"
	"alex/utils"
	"errors"
)

type ISession interface {
	GetExpire(key interface{}) (int64, error)
	SetExpire(time time.Duration) error
	SetId(id string) string
	GetName() string
	Set(key interface{}, value interface{}) error
	Get(key interface{}) (string, error)
	Delete(key interface{}) error
	Clear() error
}

const (
	SESSION_NAME = "ALEX_SESSION"
)

//Session Run on the redis
type RedisSession struct {
	id     string
	expire time.Duration
	conn   *redis.Conn
}

func NewRedisSession(c *redis.Conn) *RedisSession {
	s := new(RedisSession)
	s.SetConn(c)
	return s
}

func (r *RedisSession)getSessionKey() string {
	return utils.StringJoin("-", r.id, r.GetName())
}


//断开连接
func (r *RedisSession)Close() {
	r.conn = nil
}

//配置连接实例
func (r *RedisSession)SetConn(conn *redis.Conn) {
	r.conn = conn
}

func (r *RedisSession)SetId(id string) {
	r.id = id
}

func (r *RedisSession)GetName() string {
	return SESSION_NAME
}

func (r *RedisSession)GetId() string {
	return r.id
}

//设置session属性
func (r *RedisSession)Set(key string, value interface{}) error {
	if r.id == "" {
		return errors.New("")
	}
	_, err := r.conn.HSet(r.getSessionKey(), key, value).Result()
	if err != nil {
		return err
	}
	r.conn.Expire(key, r.expire).Result()
	return nil
}

//获取session属性
func (r *RedisSession)Get(key string) (string, error) {
	result, err := r.conn.HGet(r.getSessionKey(), key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

//获取session的所有属性
func (r *RedisSession)GetAll() (map[string]string, error) {
	result, err := r.conn.HGetAll(r.getSessionKey()).Result()
	if err != nil {
		return map[string]string{}, err
	}
	return result, nil
}

//删除session属性
func (r *RedisSession)Delete(key string) (bool, error) {
	_, err := r.conn.HDel(r.getSessionKey(), key).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

//清空用户session
func (r *RedisSession)Clear() (bool, error) {
	_, err := r.conn.Del(r.getSessionKey()).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}



