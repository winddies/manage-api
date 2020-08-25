package models

import (
	"context"
	"fmt"
	"winddies/manage-api/global"

	"github.com/go-redis/redis"
)

var RedisDb *redis.Client
var ctx = context.Background()

func RedisInit() {
	conf := global.Conf.Redis
	fmt.Println()
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	_, err := RedisDb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("connect redis success...")
	}
}
