package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func GlobalRateLimit() gin.HandlerFunc {
	rate, _ := limiter.NewRateFromFormatted("1000-M")
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	return mgin.NewMiddleware(instance)
}

func StrictRateLimit() gin.HandlerFunc {
	rate, _ := limiter.NewRateFromFormatted("1000-M")
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	return mgin.NewMiddleware(instance)
}

func SuperStrictRateLimit() gin.HandlerFunc {
	rate, _ := limiter.NewRateFromFormatted("1000-M")
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	return mgin.NewMiddleware(instance)
}