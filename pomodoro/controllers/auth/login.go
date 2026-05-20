package auth

import (
	"net/http"

	"studytrack/config"
	"studytrack/models"
	"studytrack/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Cari user berdasarkan email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Email atau password salah!",
		})
		return
	}

	// Cek akun sudah diverifikasi atau belum
	if !user.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Akun belum diverifikasi! Silakan cek email kamu.",
			"data": gin.H{
				"user_id": user.ID,
			},
		})
		return
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Email atau password salah!",
		})
		return
	}

	// Generate JWT Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal generate token!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login berhasil!",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":          user.ID,
				"nama":        user.Nama,
				"email":       user.Email,
				"universitas": user.Universitas,
				"jurusan":     user.Jurusan,
				"avatar":      user.Avatar,
			},
		},
	})
}