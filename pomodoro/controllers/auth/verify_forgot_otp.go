package auth

import (
	"fmt"
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"
	"studytrack/utils"

	"github.com/gin-gonic/gin"
)

func VerifyForgotOTP(c *gin.Context) {
	var input struct {
		UserID uint   `json:"user_id" binding:"required"`
		Kode   string `json:"kode" binding:"required"`
	}

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Cari OTP di database
	var otp models.OTP
	if err := config.DB.Where("user_id = ? AND is_used = false", input.UserID).Last(&otp).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "OTP tidak ditemukan!",
		})
		return
	}

	// Cek maksimal percobaan
	if otp.Percobaan >= 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "OTP sudah hangus! Silakan minta kirim ulang OTP.",
		})
		return
	}

	// Cek OTP expired
	if time.Now().After(otp.ExpiredAt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "OTP sudah kadaluarsa! Silakan minta kirim ulang OTP.",
		})
		return
	}

	// Cek kode OTP
	if otp.Kode != input.Kode {
		// Tambah percobaan
		config.DB.Model(&otp).Update("percobaan", otp.Percobaan+1)

		sisaPercobaan := 3 - (otp.Percobaan + 1)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("OTP salah! Sisa percobaan: %d", sisaPercobaan),
		})
		return
	}

	// OTP valid — tandai sudah dipakai
	config.DB.Model(&otp).Update("is_used", true)

	// Generate reset token
	resetToken, err := utils.GenerateResetToken(input.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal generate reset token!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "OTP valid! Silakan buat password baru.",
		"data": gin.H{
			"reset_token": resetToken,
		},
	})
}