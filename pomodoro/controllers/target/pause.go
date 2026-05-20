package target

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Pause(c *gin.Context) {
	userID, _ := c.Get("user_id")
	targetID := c.Param("id")

	// Cari target di database
	var target models.DailyTarget
	if err := config.DB.Preload("Matkul").First(&target, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Target tidak ditemukan!",
		})
		return
	}

	// Cek target milik user yang login
	if target.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Kamu tidak punya akses ke target ini!",
		})
		return
	}

	// Cek status target
	if target.Status != "running" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Target tidak bisa di-pause! Status target saat ini: " + target.Status,
		})
		return
	}

	// Catat waktu pause
	now := time.Now()

	// Update status target jadi paused
	config.DB.Model(&target).Updates(map[string]interface{}{
		"status":    "paused",
		"paused_at": now,
	})

	// Update status mata kuliah yang sedang running jadi paused
	config.DB.Model(&models.TargetMatkul{}).
		Where("target_id = ? AND status = ?", target.ID, "running").
		Update("status", "paused")

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Target di-pause! Jangan lupa lanjutkan ya! 😊",
		"data": gin.H{
			"target_id": target.ID,
			"status":    "paused",
			"paused_at": now,
		},
	})
}