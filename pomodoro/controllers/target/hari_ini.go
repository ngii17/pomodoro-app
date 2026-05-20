package target

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func HariIni(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Ambil tanggal hari ini otomatis
	hariIni := time.Now().Format("2006-01-02")

	// Cari target hari ini
	var target models.DailyTarget
	if err := config.DB.Preload("Matkul").
		Where("user_id = ? AND DATE(tanggal) = ?", userID, hariIni).
		First(&target).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Belum ada target untuk hari ini!",
			"data":    nil,
		})
		return
	}

	// Hitung progress keseluruhan
	totalSesi := 0
	totalSelesai := 0

	for _, m := range target.Matkul {
		totalSesi += m.JumlahSesi
		totalSelesai += m.SesiSelesai
	}

	// Hitung persentase progress
	progress := 0
	if totalSesi > 0 {
		progress = (totalSelesai * 100) / totalSesi
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"target": target,
			"progress": gin.H{
				"total_sesi":    totalSesi,
				"total_selesai": totalSelesai,
				"persentase":    progress,
			},
		},
	})
}