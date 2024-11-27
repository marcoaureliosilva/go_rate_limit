package limiter

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Limiter struct {
	redisClient    *redis.Client
	rateLimitIP    int
	rateLimitToken int
}

func NewLimiter() *Limiter {
	rateLimitIP, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_IP"))
	rateLimitToken, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_TOKEN"))

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})

	return &Limiter{
		redisClient:    rdb,
		rateLimitIP:    rateLimitIP,
		rateLimitToken: rateLimitToken,
	}
}

func (l *Limiter) LimitByIP(ip string) bool {
	key := "rate_limit_ip:" + ip
	return l.checkLimit(key, l.rateLimitIP)
}

func (l *Limiter) LimitByToken(token string) bool {
	key := "rate_limit_token:" + token
	return l.checkLimit(key, l.rateLimitToken)
}

func (l *Limiter) checkLimit(key string, limit int) bool {
	count, err := l.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return false
	}

	if count == 1 {
		l.redisClient.Expire(ctx, key, time.Second*1)
	}

	return count <= int64(limit)
}
