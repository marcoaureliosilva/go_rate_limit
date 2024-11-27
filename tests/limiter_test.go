package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"go_rate_limit/limiter"
)

var (
	ctx         = context.Background()
	testLimiter *limiter.Limiter
)

func setup() {
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("RATE_LIMIT_IP", "5")
	os.Setenv("RATE_LIMIT_TOKEN", "10")

	testLimiter = limiter.NewLimiter()
}

func TestLimitByIP(t *testing.T) {
	setup()

	ip := "192.168.1.1"
	for i := 0; i < 5; i++ {
		allowed := testLimiter.LimitByIP(ip)
		if !allowed {
			t.Fatalf("Expected request %d from IP %s to be allowed", i+1, ip)
		}
	}

	// The 6th request should be blocked
	allowed := testLimiter.LimitByIP(ip)
	if allowed {
		t.Fatalf("Expected the 6th request from IP %s to be blocked", ip)
	}

	// Clean up Redis
	testLimiter.redisClient.Del(ctx, "rate_limit_ip:"+ip)
}

func TestLimitByToken(t *testing.T) {
	setup()

	token := "abc123"
	for i := 0; i < 10; i++ {
		allowed := testLimiter.LimitByToken(token)
		if !allowed {
			t.Fatalf("Expected request %d with token %s to be allowed", i+1, token)
		}
	}

	// The 11th request should be blocked
	allowed := testLimiter.LimitByToken(token)
	if allowed {
		t.Fatalf("Expected the 11th request with token %s to be blocked", token)
	}

	// Clean up Redis
	testLimiter.redisClient.Del(ctx, "rate_limit_token:"+token)
}

func TestLimitWithExpiration(t *testing.T) {
	setup()

	ip := "192.168.1.2"
	token := "xyz789"

	// Allow a request from IP and Token
	allowedIP := testLimiter.LimitByIP(ip)
	allowedToken := testLimiter.LimitByToken(token)

	if !allowedIP || !allowedToken {
		t.Fatalf("Initial requests from IP %s or Token %s should be allowed", ip, token)
	}

	// Wait for 1 second before testing limits again
	time.Sleep(1 * time.Second)

	// Both limits should reset
	allowedIP = testLimiter.LimitByIP(ip)
	allowedToken = testLimiter.LimitByToken(token)

	if !allowedIP || !allowedToken {
		t.Fatalf("Requests from IP %s or Token %s should be allowed after reset", ip, token)
	}

	// Clean up Redis
	testLimiter.redisClient.Del(ctx, "rate_limit_ip:"+ip)
	testLimiter.redisClient.Del(ctx, "rate_limit_token:"+token)
}
