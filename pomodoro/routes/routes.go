package routes

import (
	authController "studytrack/controllers/auth"
	pomodoroController "studytrack/controllers/pomodoro"
	statistikController "studytrack/controllers/statistik"
	targetController "studytrack/controllers/target"
	"studytrack/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Global rate limit — semua endpoint
	r.Use(middleware.GlobalRateLimit())

	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", middleware.SuperStrictRateLimit(), authController.Register)
			auth.POST("/verify-otp", authController.VerifyOTP)
			auth.POST("/resend-otp", middleware.SuperStrictRateLimit(), authController.ResendOTP)
			auth.POST("/login", middleware.StrictRateLimit(), authController.Login)

			// Forgot Password
			auth.POST("/forgot-password", middleware.SuperStrictRateLimit(), authController.ForgotPassword)
			auth.POST("/verify-forgot-otp", authController.VerifyForgotOTP)
			auth.POST("/reset-password", authController.ResetPassword)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/test", func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				c.JSON(200, gin.H{
					"status":  "success",
					"message": "token valid",
					"user_id": userID,
				})
			})

			// Pomodoro routes
			pomodoro := protected.Group("/pomodoro")
			{
				pomodoro.POST("/mulai", pomodoroController.Mulai)
				pomodoro.PUT("/:id/pause", pomodoroController.Pause)
				pomodoro.PUT("/:id/lanjut", pomodoroController.Lanjut)
				pomodoro.PUT("/:id/selesai", pomodoroController.Selesai)
				pomodoro.PUT("/:id/skip", pomodoroController.Skip)
				pomodoro.GET("/history", pomodoroController.History)
				pomodoro.GET("/aktif", pomodoroController.Aktif)
			}

			// Target routes
			target := protected.Group("/target")
			{
				target.POST("/buat", targetController.Buat)
				target.GET("/hari-ini", targetController.HariIni)
				target.PUT("/:id/mulai", targetController.Mulai)
				target.PUT("/:id/pause", targetController.Pause)
				target.PUT("/:id/lanjut", targetController.Lanjut)
				target.GET("/history", targetController.History)
			}

			// Statistik routes
			statistik := protected.Group("/statistik")
			{
				statistik.GET("/harian", statistikController.Harian)
				statistik.GET("/mingguan", statistikController.Mingguan)
				statistik.GET("/bulanan", statistikController.Bulanan)
				statistik.GET("/matkul", statistikController.Matkul)
				statistik.GET("/streak", statistikController.Streak)
				statistik.GET("/perbandingan", statistikController.Perbandingan)
			}
		}
	}
}