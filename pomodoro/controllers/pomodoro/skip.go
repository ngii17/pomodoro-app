package pomodoro

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Skip(c *gin.Context) {
	userID, _ := c.Get("user_id")
	sesiID := c.Param("id")

	// Cari sesi di database
	var sesi models.Pomodoro
	if err := config.DB.First(&sesi, sesiID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Sesi tidak ditemukan!",
		})
		return
	}

	// Cek sesi milik user yang login
	if sesi.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Kamu tidak punya akses ke sesi ini!",
		})
		return
	}

	// Cek status sesi
	if sesi.Status != "running" && sesi.Status != "paused" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Sesi tidak bisa di-skip! Status sesi saat ini: " + sesi.Status,
		})
		return
	}

	// Catat waktu skip
	now := time.Now()

	// Update status jadi skipped
	config.DB.Model(&sesi).Updates(map[string]interface{}{
		"status":     "skipped",
		"selesai_at": now,
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sesi di-skip!",
		"data": gin.H{
			"sesi_id":     sesi.ID,
			"mata_kuliah": sesi.MataKuliah,
			"status":      "skipped",
			"selesai_at":  now,
		},
	})
}