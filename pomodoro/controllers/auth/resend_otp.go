package auth

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"
	"studytrack/utils"

	"github.com/gin-gonic/gin"
)

func ResendOTP(c *gin.Context) {
	var input struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Cari user
	var user models.User
	if err := config.DB.First(&user, input.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User tidak ditemukan!",
		})
		return
	}

	// Cek apakah akun sudah diverifikasi
	if user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Akun sudah diverifikasi!",
		})
		return
	}

	// Hanguskan semua OTP lama
	config.DB.Model(&models.OTP{}).
		Where("user_id = ? AND is_used = false", input.UserID).
		Update("is_used", true)

	// Generate OTP baru
	otpKode := utils.GenerateOTP()
	expiredAt := time.Now().Add(5 * time.Minute)

	otp := models.OTP{
		UserID:    user.ID,
		Kode:      otpKode,
		ExpiredAt: expiredAt,
		Percobaan: 0,
	}

	config.DB.Create(&otp)

	// Kirim OTP baru ke email
	if err := utils.SendOTPEmail(user.Email, user.Nama, otpKode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengirim email OTP!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "OTP baru sudah dikirim ke email kamu!",
	})
}