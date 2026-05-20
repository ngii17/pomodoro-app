package statistik

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Streak(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Ambil semua tanggal yang user pernah belajar (completed)
	var sesi []models.Pomodoro
	config.DB.Where("user_id = ? AND status = ?", userID, "completed").
		Order("created_at asc").
		Find(&sesi)

	// Kumpulkan tanggal unik yang user belajar
	tanggalBelajar := map[string]bool{}
	for _, s := range sesi {
		tanggal := s.CreatedAt.Format("2006-01-02")
		tanggalBelajar[tanggal] = true
	}

	// Hitung streak saat ini
	streakSekarang := 0
	tanggalMulaiStreak := ""
	today := time.Now()

	for i := 0; ; i++ {
		tanggal := today.AddDate(0, 0, -i).Format("2006-01-02")
		if tanggalBelajar[tanggal] {
			streakSekarang++
			tanggalMulaiStreak = tanggal
		} else {
			break
		}
	}

	// Hitung streak terpanjang
	streakTerpanjang := 0
	streakSementara := 0
	var tanggalList []string

	for tanggal := range tanggalBelajar {
		tanggalList = append(tanggalList, tanggal)
	}

	// Sort tanggal
	for i := 0; i < len(tanggalList); i++ {
		for j := i + 1; j < len(tanggalList); j++ {
			if tanggalList[i] > tanggalList[j] {
				tanggalList[i], tanggalList[j] = tanggalList[j], tanggalList[i]
			}
		}
	}

	// Hitung streak terpanjang dari list tanggal
	for i, tanggal := range tanggalList {
		if i == 0 {
			streakSementara = 1
		} else {
			prevTanggal, _ := time.Parse("2006-01-02", tanggalList[i-1])
			currTanggal, _ := time.Parse("2006-01-02", tanggal)
			selisih := int(currTanggal.Sub(prevTanggal).Hours() / 24)

			if selisih == 1 {
				streakSementara++
			} else {
				streakSementara = 1
			}
		}

		if streakSementara > streakTerpanjang {
			streakTerpanjang = streakSementara
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"streak_sekarang": gin.H{
				"hari":        streakSekarang,
				"mulai_dari":  tanggalMulaiStreak,
			},
			"streak_terpanjang": streakTerpanjang,
			"total_hari_belajar": len(tanggalBelajar),
		},
	})
}