package bootstrap

import (
	"service/app/config"
	"service/package/redis"
	"service/package/redsync"

	rsync "github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func (bs *Bootstrap) initRedis() {
	redis, redsync := getRedisClients()
	bs.Redis = redis
	bs.Redsync = redsync
}

func getRedisClients() (*redis.Client, *redsync.Redsync) {
	client := redis.Config(&redis.Setup{
		RedisHost:      config.RedisHostPath,
		RedisPassword:  config.RedisPassword,
		OpenTracingLog: !config.IsProduction(),
	})

	pool := goredis.NewPool(client.Client)
	rs := rsync.New(pool)
	return client, redsync.New(!config.IsProduction(), rs)
}
