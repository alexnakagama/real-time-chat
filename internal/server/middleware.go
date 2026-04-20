package server

import (
	"net/http"
	"sync"
	"time"
)

var (
	visitors  = make(map[string]*visitor)
	mu        sync.Mutex
	rateLimit = 10
)

type visitor struct {
	lastSeen time.Time
	requests int
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		mu.Lock()
		v, exists := visitors[ip]
		if !exists || time.Since(v.lastSeen) > time.Minute {
			visitors[ip] = &visitor{lastSeen: time.Now(), requests: 1}
		} else {
			if v.requests >= rateLimit {
				mu.Unlock()
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			v.requests++
			v.lastSeen = time.Now()
		}
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
