package statistik

import (
	"fmt"
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Bulanan(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Hitung tanggal awal & akhir bulan ini
	now := time.Now()
	awalBulan := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	akhirBulan := awalBulan.AddDate(0, 1, -1)

	awalBulanStr := awalBulan.Format("2006-01-02")
	akhirBulanStr := akhirBulan.Format("2006-01-02")

	// Ambil semua sesi bulan ini
	var sesi []models.Pomodoro
	config.DB.Where("user_id = ? AND DATE(created_at) BETWEEN ? AND ?",
		userID, awalBulanStr, akhirBulanStr).
		Find(&sesi)

	// Hitung total sesi per minggu
	sesiPerMinggu := map[string]int{}
	totalMenitBelajar := 0
	totalSesiSelesai := 0

	for _, s := range sesi {
		if s.Status == "completed" {
			// Hitung minggu ke berapa dalam bulan ini
			mingguKe := ((s.CreatedAt.Day() - 1) / 7) + 1
			key := fmt.Sprintf("Minggu %d", mingguKe)
			sesiPerMinggu[key]++
			totalMenitBelajar += s.DurasiBelajar
			totalSesiSelesai++
		}
	}

	// Tentukan minggu paling produktif
	mingguTerproduktif := ""
	maxSesi := 0
	for minggu, jumlah := range sesiPerMinggu {
		if jumlah > maxSesi {
			maxSesi = jumlah
			mingguTerproduktif = minggu
		}
	}

	// Hitung rata-rata sesi per hari
	jumlahHari := int(akhirBulan.Sub(awalBulan).Hours()/24) + 1
	rataRata := 0
	if jumlahHari > 0 {
		rataRata = totalSesiSelesai / jumlahHari
	}

	// Konversi menit ke jam & menit
	jam := totalMenitBelajar / 60
	menit := totalMenitBelajar % 60

	// Hitung total target selesai bulan ini
	var totalTargetSelesai int64
	config.DB.Model(&models.DailyTarget{}).
		Where("user_id = ? AND DATE(tanggal) BETWEEN ? AND ? AND status = ?",
			userID, awalBulanStr, akhirBulanStr, "completed").
		Count(&totalTargetSelesai)

	// Buat data per minggu untuk grafik
	var grafikMingguan []gin.H
	for i := 1; i <= 5; i++ {
		key := fmt.Sprintf("Minggu %d", i)
		grafikMingguan = append(grafikMingguan, gin.H{
			"minggu": key,
			"sesi":   sesiPerMinggu[key],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"periode": gin.H{
				"bulan": now.Format("January 2006"),
				"awal":  awalBulanStr,
				"akhir": akhirBulanStr,
			},
			"total_sesi_selesai":      totalSesiSelesai,
			"total_belajar": gin.H{
				"jam":   jam,
				"menit": menit,
			},
			"rata_rata_sesi_per_hari": rataRata,
			"minggu_terproduktif":     mingguTerproduktif,
			"total_target_selesai":    totalTargetSelesai,
			"grafik_mingguan":         grafikMingguan,
		},
	})
}