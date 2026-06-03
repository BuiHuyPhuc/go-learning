package middlewares

import (
	"fmt"
	"go-learning/global"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"

	redisStore "github.com/ulule/limiter/v3/drivers/store/redis"
)

type RateLimiter struct {
	globalRateLimiter         *limiter.Limiter
	publicAPIRateLimiter      *limiter.Limiter
	userPrivateAPIRateLimiter *limiter.Limiter
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		globalRateLimiter:         rateLimiter("100-S"),
		publicAPIRateLimiter:      rateLimiter("80-S"),
		userPrivateAPIRateLimiter: rateLimiter("50-S"),
	}
}

func rateLimiter(interval string) *limiter.Limiter {
	store, err := redisStore.NewStoreWithOptions(global.Rdb, limiter.StoreOptions{
		Prefix:          "rate-limiter", // u:uuid -> u:1001
		MaxRetry:        3,
		CleanUpInterval: time.Hour,
	})
	if err != nil {
		return nil
	}

	rate, err := limiter.NewRateFromFormatted(interval) // 5-S, 10-M
	if err != nil {
		return nil
	}

	instance := limiter.New(store, rate)
	return instance
}

// Global API limiter
func (rl *RateLimiter) GlobalRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "global"
		log.Println("global --->")
		limitContext, err := rl.globalRateLimiter.Get(c, key)
		if err != nil {
			fmt.Println("Failed to check rate limit GLOBAL: ", err)
			c.Next()
			return
		}

		if limitContext.Reached {
			log.Println("Rate limit breached GLOBAL")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit breached GLOBAL, try later"})
		}

		c.Next()
	}
}

// Public API limiter
func (rl *RateLimiter) PublicRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPath := c.Request.URL.Path // urlPath: /ping/80
		rateLimitPath := rl.filterLimitUrlPath(urlPath)
		if rateLimitPath != nil {
			log.Println("Client Ip --->", c.ClientIP())
			key := fmt.Sprintf("%s-%s", "111-222-333-44", urlPath)
			limitContext, err := rateLimitPath.Get(c, key)
			if err != nil {
				fmt.Println("Failed to check rate limit PUBLIC: ", err)
				c.Next()
				return
			}

			if limitContext.Reached {
				log.Println("Rate limit breached PUBLIC")
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit breached PUBLIC, try later"})
			}
		}

		c.Next()
	}
}

// User private API limiter
func (rl *RateLimiter) UserPrivateRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPath := c.Request.URL.Path // urlPath: /ping/50
		rateLimitPath := rl.filterLimitUrlPath(urlPath)
		if rateLimitPath != nil {
			userId := 1001 // context.GetUserIdFromUUID(c.Request.Context())
			key := fmt.Sprintf("%d-%s", userId, urlPath)
			limitContext, err := rateLimitPath.Get(c, key)
			if err != nil {
				fmt.Println("Failed to check rate limit USER PRIVATE: ", err)
				c.Next()
				return
			}

			if limitContext.Reached {
				log.Println("Rate limit breached USER PRIVATE")
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit breached USER PRIVATE, try later"})
			}
		}

		c.Next()
	}
}

func (rl *RateLimiter) filterLimitUrlPath(urlPath string) *limiter.Limiter {
	if urlPath == "/api/v1/users/login" || urlPath == "/ping/80" {
		return rl.publicAPIRateLimiter
	} else if urlPath == "/api/v1/users/get-info" || urlPath == "/ping/50" {
		return rl.userPrivateAPIRateLimiter
	} else {
		return rl.globalRateLimiter
	}
}
