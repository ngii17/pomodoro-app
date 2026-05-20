package auth

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"
	"studytrack/utils"

	"github.com/gin-gonic/gin"
)

func ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Cek email ada di database
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Email tidak ditemukan!",
		})
		return
	}

	// Hanguskan OTP lama
	config.DB.Model(&models.OTP{}).
		Where("user_id = ? AND is_used = false", user.ID).
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

	// Kirim OTP ke email
	if err := utils.SendForgotPasswordEmail(user.Email, user.Nama, otpKode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengirim email OTP!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "OTP reset password sudah dikirim ke email kamu!",
		"data": gin.H{
			"user_id": user.ID,
		},
	})
}