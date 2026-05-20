package pomodoro

import (
	"net/http"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Pause(c *gin.Context) {
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

	// Cek status sesi — hanya bisa pause kalau running
	if sesi.Status != "running" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Sesi tidak bisa di-pause! Status sesi saat ini: " + sesi.Status,
		})
		return
	}

	// Update status & tambah total pause
	config.DB.Model(&sesi).Updates(map[string]interface{}{
		"status":      "paused",
		"total_pause": sesi.TotalPause + 1,
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sesi di-pause!",
		"data": gin.H{
			"sesi_id":     sesi.ID,
			"status":      "paused",
			"total_pause": sesi.TotalPause + 1,
		},
	})
}