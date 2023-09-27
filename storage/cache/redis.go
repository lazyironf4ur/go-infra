package cache

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/lazyironf4ur/go-infra/common"
	"github.com/lazyironf4ur/go-infra/conf"
)

var RedisClient *redis.Client

func init() {
	if conf.GlobalConfig["redis"] != nil {
		redisConf := conf.GlobalConfig["redis"].(map[string]interface{})
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisConf["host"], redisConf["port"]),
			Password: fmt.Sprintf("%s", redisConf["password"]),
			DB:       redisConf["db"].(int),
		})
		common.Must(RedisClient)
	}
}
