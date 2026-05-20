package middleware

import (
	"net/http"
	"strings"

	"studytrack/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header
		authHeader := c.GetHeader("Authorization")

		// Cek header ada atau tidak
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token tidak ditemukan! Silakan login.",
			})
			c.Abort()
			return
		}

		// Format header harus "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Format token tidak valid!",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validasi token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token tidak valid atau sudah expired! Silakan login ulang.",
			})
			c.Abort()
			return
		}

		// Simpan user_id ke context
		c.Set("user_id", claims.UserID)

		// Lanjut ke controller
		c.Next()
	}
}