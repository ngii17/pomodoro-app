package auth

import (
	"net/http"

	"studytrack/config"
	"studytrack/models"
	"studytrack/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ResetPassword(c *gin.Context) {
	var input struct {
		ResetToken      string `json:"reset_token" binding:"required"`
		Password        string `json:"password" binding:"required,min=8"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Cek password & konfirmasi password cocok
	if input.Password != input.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Password dan konfirmasi password tidak cocok!",
		})
		return
	}

	// Validasi reset token
	claims, err := utils.ValidateToken(input.ResetToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Reset token tidak valid atau sudah expired! Silakan ulangi proses forgot password.",
		})
		return
	}

	// Cari user di database
	var user models.User
	if err := config.DB.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User tidak ditemukan!",
		})
		return
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memproses password!",
		})
		return
	}

	// Update password di database
	if err := config.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal update password!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password berhasil direset! Silakan login dengan password baru.",
	})
}