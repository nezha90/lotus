package main

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

var (
	// 每个IP地址的限速器
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

// 获取或创建限速器
func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(1, 5) // 每秒1个请求，最多允许5个突发请求
		limiters[ip] = limiter
	}
	return limiter
}

// 限速中间件
func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
