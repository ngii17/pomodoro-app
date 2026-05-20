package pomodoro

import (
	"net/http"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Aktif(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Cari sesi yang masih aktif
	var sesi models.Pomodoro
	if err := config.DB.Where("user_id = ? AND status IN ?", userID, []string{"running", "paused"}).
		First(&sesi).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Tidak ada sesi yang sedang aktif!",
			"data":    nil,
		})
		return
	}

	// Hitung sisa waktu
	var pesanStatus string
	if sesi.Status == "running" {
		pesanStatus = "Sesi sedang berjalan! Semangat! 💪"
	} else {
		pesanStatus = "Sesi sedang di-pause!"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": pesanStatus,
		"data": gin.H{
			"sesi_id":                sesi.ID,
			"mata_kuliah":            sesi.MataKuliah,
			"metode":                 sesi.Metode,
			"status":                 sesi.Status,
			"durasi_belajar":         sesi.DurasiBelajar,
			"durasi_istirahat":       sesi.DurasiIstirahat,
			"durasi_istirahat_panjang": sesi.DurasiIstirahatPanjang,
			"sesi_keberapa":          sesi.SesiKeberapa,
			"jumlah_sesi":            sesi.JumlahSesi,
			"total_pause":            sesi.TotalPause,
			"mulai_at":               sesi.MulaiAt,
		},
	})
}