package target

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Lanjut(c *gin.Context) {
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
	if target.Status != "paused" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Target tidak bisa dilanjutkan! Status target saat ini: " + target.Status,
		})
		return
	}

	// Hitung durasi pause
	now := time.Now()
	pauseDuration := 0
	if target.PausedAt != nil {
		pauseDuration = int(now.Sub(*target.PausedAt).Minutes())
	}

	// Update status target jadi running
	config.DB.Model(&target).Updates(map[string]interface{}{
		"status":               "running",
		"paused_at":            nil,
		"total_pause_duration": target.TotalPauseDuration + pauseDuration,
	})

	// Update status mata kuliah yang paused jadi running kembali
	var matkulAktif models.TargetMatkul
	config.DB.Where("target_id = ? AND status = ?", target.ID, "paused").
		First(&matkulAktif)
	config.DB.Model(&matkulAktif).Update("status", "running")

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Target dilanjutkan! Semangat belajar! 💪",
		"data": gin.H{
			"target_id":            target.ID,
			"status":               "running",
			"total_pause_duration": target.TotalPauseDuration + pauseDuration,
			"matkul_sekarang": gin.H{
				"id":               matkulAktif.ID,
				"mata_kuliah":      matkulAktif.MataKuliah,
				"jumlah_sesi":      matkulAktif.JumlahSesi,
				"sesi_selesai":     matkulAktif.SesiSelesai,
				"durasi_istirahat": matkulAktif.DurasiIstirahat,
				"urutan":           matkulAktif.Urutan,
			},
		},
	})
}