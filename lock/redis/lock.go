package utils

import (
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"log"
	"time"
)


// 业务锁 使用redis 来设置
type RedisLocker struct {
	Redis *redis.Client
	Token string //共享token， 每台机器的token 不一样，那么删除的时候先判断这个值，哪台机器才有删除的权限
}

func MustLocker(locker *RedisLocker, err error) *RedisLocker {
	if err != nil {
		panic(err)
	}
	return locker
}

func NewLocker(client *redis.Client) (*RedisLocker, error) {

	return &RedisLocker{
		client,
		uuid.New().String()}, nil
}

func (l *RedisLocker) TryLock(key string, expire time.Duration) bool {
	res, err := l.Redis.SetNX(key, l.Token, expire).Result()
	if err != nil {
		log.Printf("lock fail %s", err.Error())
		return false
	}
	return res
}

func (l *RedisLocker) Release(key string) bool {
	luaScript := "if redis.call('get',KEYS[1]) == ARGV[1] then " +
		"return redis.call('del',KEYS[1]) else return 0 end"
	_, err := l.Redis.Eval(luaScript, []string{key}, l.Token).Result()
	if err != nil {
		log.Printf("release lock fail %s", err.Error())
		return false
	}
	return true
	//script := redis.NewScript(1, luaScript)
	//conn := l.conn()
	//if conn == nil {
	//	return false
	//}
	//defer func() {
	//	_ = conn.Close()
	//}()
	//if code, err := redis.Int(script.Do(conn, key, l.Token)); err != nil {
	//	return false
	//} else {
	//	return code == 1
	//}
}
