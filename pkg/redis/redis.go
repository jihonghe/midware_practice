package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	Client *redis.Client
)

func Init(config *viper.Viper) {
	var opt = &redis.Options{
		Addr:     config.GetString("redis.addr"),
		Password: config.GetString("redis.password"),
	}
	Client = redis.NewClient(opt)
}
