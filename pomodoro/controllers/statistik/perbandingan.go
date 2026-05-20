package statistik

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Perbandingan(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Hitung periode minggu ini
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	awalMingguIni := now.AddDate(0, 0, -(weekday - 1))
	akhirMingguIni := now.AddDate(0, 0, 7-weekday)

	// Hitung periode minggu lalu
	awalMingguLalu := awalMingguIni.AddDate(0, 0, -7)
	akhirMingguLalu := akhirMingguIni.AddDate(0, 0, -7)

	// Format tanggal
	awalMingguIniStr := awalMingguIni.Format("2006-01-02")
	akhirMingguIniStr := akhirMingguIni.Format("2006-01-02")
	awalMingguLaluStr := awalMingguLalu.Format("2006-01-02")
	akhirMingguLaluStr := akhirMingguLalu.Format("2006-01-02")

	// Fungsi helper hitung statistik per periode
	hitungStatistik := func(awal, akhir string) (int, int, string) {
		var sesi []models.Pomodoro
		config.DB.Where("user_id = ? AND DATE(created_at) BETWEEN ? AND ? AND status = ?",
			userID, awal, akhir, "completed").
			Find(&sesi)

		totalSesi := 0
		totalMenit := 0
		hariNames := []string{"Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Minggu"}
		sesiPerHari := map[string]int{}

		for _, s := range sesi {
			totalSesi++
			totalMenit += s.DurasiBelajar
			hariIndex := int(s.CreatedAt.Weekday())
			if hariIndex == 0 {
				hariIndex = 7
			}
			sesiPerHari[hariNames[hariIndex-1]]++
		}

		hariTerproduktif := ""
		maxSesi := 0
		for hari, jumlah := range sesiPerHari {
			if jumlah > maxSesi {
				maxSesi = jumlah
				hariTerproduktif = hari
			}
		}

		return totalSesi, totalMenit, hariTerproduktif
	}

	// Hitung statistik minggu ini & minggu lalu
	sesiMingguIni, menitMingguIni, hariMingguIni := hitungStatistik(awalMingguIniStr, akhirMingguIniStr)
	sesiMingguLalu, menitMingguLalu, hariMingguLalu := hitungStatistik(awalMingguLaluStr, akhirMingguLaluStr)

	// Hitung persentase naik/turun
	persentase := 0
	tren := "sama"
	if sesiMingguLalu > 0 {
		persentase = ((sesiMingguIni - sesiMingguLalu) * 100) / sesiMingguLalu
		if persentase > 0 {
			tren = "naik"
		} else if persentase < 0 {
			tren = "turun"
			persentase = -persentase
		}
	} else if sesiMingguIni > 0 {
		tren = "naik"
		persentase = 100
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"minggu_ini": gin.H{
				"periode": gin.H{
					"awal":  awalMingguIniStr,
					"akhir": akhirMingguIniStr,
				},
				"total_sesi": sesiMingguIni,
				"total_belajar": gin.H{
					"jam":   menitMingguIni / 60,
					"menit": menitMingguIni % 60,
				},
				"hari_terproduktif": hariMingguIni,
			},
			"minggu_lalu": gin.H{
				"periode": gin.H{
					"awal":  awalMingguLaluStr,
					"akhir": akhirMingguLaluStr,
				},
				"total_sesi": sesiMingguLalu,
				"total_belajar": gin.H{
					"jam":   menitMingguLalu / 60,
					"menit": menitMingguLalu % 60,
				},
				"hari_terproduktif": hariMingguLalu,
			},
			"perbandingan": gin.H{
				"tren":       tren,
				"persentase": persentase,
			},
		},
	})
}