package statistik

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Mingguan(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Hitung tanggal awal & akhir minggu ini (Senin - Minggu)
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	awalMinggu := now.AddDate(0, 0, -(weekday - 1)).Format("2006-01-02")
	akhirMinggu := now.AddDate(0, 0, 7-weekday).Format("2006-01-02")

	// Ambil semua sesi minggu ini
	var sesi []models.Pomodoro
	config.DB.Where("user_id = ? AND DATE(created_at) BETWEEN ? AND ?", userID, awalMinggu, akhirMinggu).
		Find(&sesi)

	// Hitung total sesi per hari
	hariNames := []string{"Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Minggu"}
	sesiPerHari := map[string]int{}
	for _, h := range hariNames {
		sesiPerHari[h] = 0
	}

	totalMenitBelajar := 0
	totalSesiSelesai := 0

	for _, s := range sesi {
		if s.Status == "completed" {
			hariIndex := int(s.CreatedAt.Weekday())
			if hariIndex == 0 {
				hariIndex = 7
			}
			hariNama := hariNames[hariIndex-1]
			sesiPerHari[hariNama]++
			totalMenitBelajar += s.DurasiBelajar
			totalSesiSelesai++
		}
	}

	// Tentukan hari paling produktif
	hariTerproduktif := ""
	maxSesi := 0
	for hari, jumlah := range sesiPerHari {
		if jumlah > maxSesi {
			maxSesi = jumlah
			hariTerproduktif = hari
		}
	}

	// Hitung rata-rata sesi per hari
	rataRata := 0
	if totalSesiSelesai > 0 {
		rataRata = totalSesiSelesai / 7
	}

	// Konversi menit ke jam & menit
	jam := totalMenitBelajar / 60
	menit := totalMenitBelajar % 60

	// Hitung total target selesai minggu ini
	var totalTargetSelesai int64
	config.DB.Model(&models.DailyTarget{}).
		Where("user_id = ? AND DATE(tanggal) BETWEEN ? AND ? AND status = ?",
			userID, awalMinggu, akhirMinggu, "completed").
		Count(&totalTargetSelesai)

	// Buat data per hari untuk grafik
	var grafikHarian []gin.H
	for _, hari := range hariNames {
		grafikHarian = append(grafikHarian, gin.H{
			"hari":  hari,
			"sesi":  sesiPerHari[hari],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"periode": gin.H{
				"awal":  awalMinggu,
				"akhir": akhirMinggu,
			},
			"total_sesi_selesai": totalSesiSelesai,
			"total_belajar": gin.H{
				"jam":   jam,
				"menit": menit,
			},
			"rata_rata_sesi_per_hari": rataRata,
			"hari_terproduktif":       hariTerproduktif,
			"total_target_selesai":    totalTargetSelesai,
			"grafik_harian":           grafikHarian,
		},
	})
}