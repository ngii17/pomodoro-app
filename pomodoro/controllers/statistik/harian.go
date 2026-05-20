package statistik

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Harian(c *gin.Context) {
	userID, _ := c.Get("user_id")
	hariIni := time.Now().Format("2006-01-02")

	// Ambil semua sesi pomodoro hari ini
	var sesi []models.Pomodoro
	config.DB.Where("user_id = ? AND DATE(created_at) = ?", userID, hariIni).
		Find(&sesi)

	// Hitung statistik sesi
	totalSesiSelesai := 0
	totalMenitBelajar := 0
	totalPause := 0
	matkulHariIni := map[string]int{}

	for _, s := range sesi {
		if s.Status == "completed" {
			totalSesiSelesai++
			totalMenitBelajar += s.DurasiBelajar
			matkulHariIni[s.MataKuliah]++
		}
		totalPause += s.TotalPause
	}

	// Konversi menit ke jam & menit
	jam := totalMenitBelajar / 60
	menit := totalMenitBelajar % 60

	// Buat daftar mata kuliah
	var daftarMatkul []gin.H
	for matkul, jumlahSesi := range matkulHariIni {
		daftarMatkul = append(daftarMatkul, gin.H{
			"mata_kuliah":  matkul,
			"jumlah_sesi": jumlahSesi,
		})
	}

	// Cek target hari ini
	var target models.DailyTarget
	targetAda := false
	targetSelesai := false

	if err := config.DB.Preload("Matkul").
		Where("user_id = ? AND DATE(tanggal) = ?", userID, hariIni).
		First(&target).Error; err == nil {
		targetAda = true
		targetSelesai = target.Status == "completed"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"tanggal": hariIni,
			"sesi": gin.H{
				"total_sesi_selesai": totalSesiSelesai,
				"total_belajar": gin.H{
					"jam":   jam,
					"menit": menit,
				},
				"total_pause": totalPause,
			},
			"mata_kuliah": daftarMatkul,
			"target": gin.H{
				"ada":     targetAda,
				"selesai": targetSelesai,
			},
		},
	})
}