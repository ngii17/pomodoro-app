package auth

import (
	"fmt"
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func VerifyOTP(c *gin.Context) {
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

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("OTP salah! Sisa percobaan: %d", 3-(otp.Percobaan+1)),
		})
		return
	}

	// OTP valid — aktifkan akun
	now := time.Now()
	config.DB.Model(&otp).Update("is_used", true)
	config.DB.Model(&models.User{}).Where("id = ?", input.UserID).Updates(map[string]interface{}{
		"is_verified": true,
		"verified_at": now,
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Email berhasil diverifikasi! Silakan login.",
	})
}