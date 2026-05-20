package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// Global rate limit — 60 request per menit
func GlobalRateLimit() gin.HandlerFunc {
	rate, _ := limiter.NewRateFromFormatted("60-M")
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	return mgin.NewMiddleware(instance)
}

// Strict rate limit — 5 request per menit (untuk login)
func StrictRateLimit() gin.HandlerFunc {
	rate, _ := limiter.NewRateFromFormatted("5-M")
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	return mgin.NewMiddleware(instance)
}

// Super strict rate limit — 3 request per jam (untuk register, forgot password, resend otp)
func SuperStrictRateLimit() gin.HandlerFunc {
	rate, _ := limiter.NewRateFromFormatted("3-H")
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	return mgin.NewMiddleware(instance)
}