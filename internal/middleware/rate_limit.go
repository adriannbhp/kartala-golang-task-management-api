package middleware

import (
	"log"
	"net/http"
	"sync"

	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		i.mu.Lock()
		defer i.mu.Unlock()

		// Re-check after acquiring lock
		limiter, exists = i.ips[ip]
		if !exists {
			limiter = rate.NewLimiter(i.r, i.b)
			i.ips[ip] = limiter
		}
	}

	return limiter
}

func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(r, b)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)
		
		// Debug log
		tokens := l.Tokens()
		log.Printf("[RateLimit] IP: %s | Tokens: %.2f", ip, tokens)

		if !l.Allow() {
			log.Printf("[RateLimit] IP %s blocked", ip)
			response.AbortWithError(c, http.StatusTooManyRequests, "Too many requests. Please try again later.")
			return
		}
		c.Next()
	}
}
