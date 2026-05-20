package target

import (
	"net/http"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func History(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Ambil query filter
	status := c.Query("status")

	// Build query
	query := config.DB.Preload("Matkul").Where("user_id = ?", userID)

	// Filter status
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Ambil semua target
	var targets []models.DailyTarget
	if err := query.Order("tanggal desc").Find(&targets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil history target!",
		})
		return
	}

	// Hitung statistik per target
	type TargetWithProgress struct {
		Target   models.DailyTarget `json:"target"`
		Progress gin.H              `json:"progress"`
	}

	var hasil []TargetWithProgress

	for _, t := range targets {
		totalSesi := 0
		totalSelesai := 0

		for _, m := range t.Matkul {
			totalSesi += m.JumlahSesi
			totalSelesai += m.SesiSelesai
		}

		persentase := 0
		if totalSesi > 0 {
			persentase = (totalSelesai * 100) / totalSesi
		}

		hasil = append(hasil, TargetWithProgress{
			Target: t,
			Progress: gin.H{
				"total_sesi":    totalSesi,
				"total_selesai": totalSelesai,
				"persentase":    persentase,
			},
		})
	}

	// Hitung ringkasan keseluruhan
	totalTarget := len(targets)
	totalCompleted := 0
	for _, t := range targets {
		if t.Status == "completed" {
			totalCompleted++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"targets": hasil,
			"ringkasan": gin.H{
				"total_target":    totalTarget,
				"total_completed": totalCompleted,
				"total_pending":   totalTarget - totalCompleted,
			},
		},
	})
}