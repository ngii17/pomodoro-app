package statistik

import (
	"net/http"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Matkul(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Ambil semua sesi yang completed
	var sesi []models.Pomodoro
	config.DB.Where("user_id = ? AND status = ?", userID, "completed").
		Find(&sesi)

	// Kelompokkan per mata kuliah
	type StatMatkul struct {
		TotalSesi    int
		TotalMenit   int
		TotalPause   int
	}

	statPerMatkul := map[string]*StatMatkul{}

	for _, s := range sesi {
		if _, exists := statPerMatkul[s.MataKuliah]; !exists {
			statPerMatkul[s.MataKuliah] = &StatMatkul{}
		}
		statPerMatkul[s.MataKuliah].TotalSesi++
		statPerMatkul[s.MataKuliah].TotalMenit += s.DurasiBelajar
		statPerMatkul[s.MataKuliah].TotalPause += s.TotalPause
	}

	// Buat daftar statistik per mata kuliah
	var daftarMatkul []gin.H
	matkulTerbanyak := ""
	maxSesi := 0

	for matkul, stat := range statPerMatkul {
		jam := stat.TotalMenit / 60
		menit := stat.TotalMenit % 60

		rataRataPause := 0
		if stat.TotalSesi > 0 {
			rataRataPause = stat.TotalPause / stat.TotalSesi
		}

		daftarMatkul = append(daftarMatkul, gin.H{
			"mata_kuliah":     matkul,
			"total_sesi":      stat.TotalSesi,
			"total_belajar": gin.H{
				"jam":   jam,
				"menit": menit,
			},
			"rata_rata_pause": rataRataPause,
		})

		if stat.TotalSesi > maxSesi {
			maxSesi = stat.TotalSesi
			matkulTerbanyak = matkul
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"matkul_terbanyak": matkulTerbanyak,
			"total_matkul":     len(daftarMatkul),
			"daftar_matkul":    daftarMatkul,
		},
	})
}