package bootstrap

import (
	"service/app/repository"
	"service/package/redis"
	"service/package/redsync"

	"github.com/RevenueMonster/sqlike/sqlike"
)

type Bootstrap struct {
	Database   *sqlike.Database
	Redsync    *redsync.Redsync
	Redis      *redis.Client
	Repository *repository.Repository
}

// New :
func New() *Bootstrap {
	bs := new(Bootstrap)
	bs.initMySQL()
	bs.initRedis()
	bs.initRepository()
	return bs
}
