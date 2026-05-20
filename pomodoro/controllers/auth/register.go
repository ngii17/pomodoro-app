package auth

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"
	"studytrack/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Nama        string `json:"nama" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=8"`
		Universitas string `json:"universitas" binding:"required"`
		Jurusan     string `json:"jurusan" binding:"required"`
		Avatar      string `json:"avatar" binding:"required"`
	}

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Cek email sudah dipakai atau belum
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Email sudah digunakan!",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memproses password!",
		})
		return
	}

	// Simpan user ke database
	user := models.User{
		Nama:        input.Nama,
		Email:       input.Email,
		Password:    string(hashedPassword),
		Universitas: input.Universitas,
		Jurusan:     input.Jurusan,
		Avatar:      input.Avatar,
		IsVerified:  false,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat akun!",
		})
		return
	}

	// Generate OTP
	otpKode := utils.GenerateOTP()
	expiredAt := time.Now().Add(5 * time.Minute)

	otp := models.OTP{
		UserID:    user.ID,
		Kode:      otpKode,
		ExpiredAt: expiredAt,
	}

	config.DB.Create(&otp)

	// Kirim OTP ke email
	if err := utils.SendOTPEmail(user.Email, user.Nama, otpKode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengirim email OTP!",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Akun berhasil dibuat! Silakan cek email kamu untuk verifikasi OTP.",
		"data": gin.H{
			"user_id": user.ID,
			"email":   user.Email,
		},
	})
}