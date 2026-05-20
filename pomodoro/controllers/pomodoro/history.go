package pomodoro

import (
	"net/http"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func History(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Ambil query filter
	tanggal    := c.Query("tanggal")     // format: 2026-05-17
	mataKuliah := c.Query("mata_kuliah")

	// Build query
	query := config.DB.Where("user_id = ?", userID)

	// Filter tanggal
	if tanggal != "" {
		query = query.Where("DATE(created_at) = ?", tanggal)
	}

	// Filter mata kuliah
	if mataKuliah != "" {
		query = query.Where("mata_kuliah = ?", mataKuliah)
	}

	// Ambil semua sesi
	var sesi []models.Pomodoro
	if err := query.Order("created_at desc").Find(&sesi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil history sesi!",
		})
		return
	}

	// Hitung total waktu belajar & statistik
	totalMenit := 0
	totalSelesai := 0
	totalSkip := 0

	for _, s := range sesi {
		if s.Status == "completed" {
			totalMenit += s.DurasiBelajar
			totalSelesai++
		} else if s.Status == "skipped" {
			totalSkip++
		}
	}

	// Konversi menit ke jam & menit
	jam := totalMenit / 60
	menit := totalMenit % 60

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"sesi": sesi,
			"ringkasan": gin.H{
				"total_sesi":    len(sesi),
				"total_selesai": totalSelesai,
				"total_skip":    totalSkip,
				"total_belajar": gin.H{
					"jam":   jam,
					"menit": menit,
				},
			},
		},
	})
}