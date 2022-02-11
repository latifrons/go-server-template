package cache

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/latifrons/lbserver/tools"
	"time"
)

type RedisCache struct {
	Address  string
	Password string
	Db       int
	Rdb      *redis.Client
}

func (r *RedisCache) Init() {
	r.Rdb = redis.NewClient(&redis.Options{
		Addr:     r.Address,
		Password: r.Password,
		DB:       r.Db,
	})
}

func (r *RedisCache) BenchmarkWrite() {
	expire := time.Second * 300
	ctx := tools.GetContext(20)
	for i := 0; i < 100000; i++ {
		err := r.Rdb.Set(ctx, fmt.Sprintf("key-%d", i), i, expire).Err()
		if err != nil {
			panic(err)
		}
	}
}

func (r *RedisCache) BenchmarkRead() {
	ctx := tools.GetContext(20)
	for i := 0; i < 100000; i++ {
		_ = r.Rdb.Get(ctx, fmt.Sprintf("key-%d", i)).Val()
		//fmt.Println("key", val)
	}
}

// GetExclusiveLock will try to get a redis lock and release it after some time
// it is necessary to get a lock and lock it for a while to globally schedule a cron.
func (r *RedisCache) GetExclusiveLock(name string, lockTime time.Duration) bool {
	nx := r.Rdb.SetNX(tools.GetContextDefault(), "lock-"+name, 1, lockTime)
	if nx.Err() != nil {
		return false
	}
	return nx.Val()
}
