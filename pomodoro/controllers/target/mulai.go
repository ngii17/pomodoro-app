package target

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Mulai(c *gin.Context) {
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
	if target.Status != "not_started" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Target tidak bisa dimulai! Status target saat ini: " + target.Status,
		})
		return
	}

	// Catat waktu mulai
	now := time.Now()

	// Update status target jadi running
	config.DB.Model(&target).Updates(map[string]interface{}{
		"status":   "running",
		"mulai_at": now,
	})

	// Update status mata kuliah pertama jadi running
	var matkulPertama models.TargetMatkul
	config.DB.Where("target_id = ? AND urutan = 1", target.ID).First(&matkulPertama)
	config.DB.Model(&matkulPertama).Update("status", "running")

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Target dimulai! Semangat belajar! 💪",
		"data": gin.H{
			"target_id": target.ID,
			"status":    "running",
			"mulai_at":  now,
			"matkul_sekarang": gin.H{
				"id":              matkulPertama.ID,
				"mata_kuliah":     matkulPertama.MataKuliah,
				"jumlah_sesi":     matkulPertama.JumlahSesi,
				"durasi_istirahat": matkulPertama.DurasiIstirahat,
				"urutan":          matkulPertama.Urutan,
			},
		},
	})
}