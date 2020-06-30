package gredis

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"reflect"
	"time"
	"userSystem/pkg/setting"
)

var RedisClient *redis.Client

func Setup() {
	RedisClient = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     setting.RedisSetting.Host,
		Password: setting.RedisSetting.Password,
	})
	val, err := RedisClient.Ping().Result()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("[info] Redis PING : %s", val)
}

//设置一个key的值
func Set(data map[string]string, second time.Duration) error {
	for k, v := range data {
		if err := RedisClient.Set(k, v, second*time.Second).Err(); err != nil {
			return err
		}
	}
	return nil
}

//针对一个key的数值进行递增操作
func Incr(key string, second time.Duration) error {
	_, err := RedisClient.Incr(key).Result()
	RedisClient.Expire(key, second)
	if err != nil {
		return err
	}
	return nil
}

//查询key的值
func Get(key string) (string, error) {
	isExists, err := Exists(key)
	if err != nil {
		return "", err
	}
	if !isExists {
		return "", nil
	}
	return RedisClient.Get(key).Result()
}

//查询key的有效期
func GetTTL(key string) (float64, error) {
	isExists, err := Exists(key)
	if err != nil {
		return 0, err
	}
	if !isExists {
		return 0, nil
	}
	times, err := RedisClient.TTL(key).Result()
	if err != nil {
		return 0, err
	}
	return times.Seconds(), err
}

//检查key是否存在
func Exists(key string) (bool, error) {
	ok, err := RedisClient.Exists(key).Result()
	if err != nil {
		return false, err
	}
	if ok == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

//删除key的值
func Delete(key string) {
	RedisClient.Del(key)
}

//模糊删除
func LikeDeletes(key string) {
	keys, err := RedisClient.Do("KEYS", "*"+key+"*").Result()
	if err != nil {
		return
	}
	if reflect.TypeOf(keys).Kind() == reflect.Slice {
		s := reflect.ValueOf(keys)
		for i := 0; i < s.Len(); i++ {
			Delete(fmt.Sprintf("%s", s.Index(i)))
		}
	}
}
