package middleware

import (
	"go_rate_limit/limiter"
	"net/http"
)

func RateLimiterMiddleware(next http.Handler) http.Handler {
	l := limiter.NewLimiter()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		token := r.Header.Get("API_KEY")

		ipAllowed := l.LimitByIP(ip)
		tokenAllowed := l.LimitByToken(token)

		if !ipAllowed && !tokenAllowed {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
