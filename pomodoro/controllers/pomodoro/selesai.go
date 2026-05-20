package pomodoro

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Selesai(c *gin.Context) {
	userID, _ := c.Get("user_id")
	sesiID := c.Param("id")

	// Cari sesi di database
	var sesi models.Pomodoro
	if err := config.DB.First(&sesi, sesiID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Sesi tidak ditemukan!",
		})
		return
	}

	// Cek sesi milik user yang login
	if sesi.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Kamu tidak punya akses ke sesi ini!",
		})
		return
	}

	// Cek status sesi
	if sesi.Status != "running" && sesi.Status != "paused" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Sesi sudah selesai atau di-skip!",
		})
		return
	}

	// Catat waktu selesai
	now := time.Now()

	// Update status jadi completed
	config.DB.Model(&sesi).Updates(map[string]interface{}{
		"status":     "completed",
		"selesai_at": now,
	})

	// Tentukan istirahat selanjutnya
	var pesanIstirahat string
	var durasiIstirahat int

	if sesi.SesiKeberapa >= sesi.JumlahSesi {
		// Sudah mencapai jumlah sesi → istirahat panjang
		pesanIstirahat = "Kerja keras! Waktunya istirahat panjang! 🛋️"
		durasiIstirahat = sesi.DurasiIstirahatPanjang
	} else {
		// Belum mencapai jumlah sesi → istirahat pendek
		pesanIstirahat = "Bagus! Waktunya istirahat sebentar! ☕"
		durasiIstirahat = sesi.DurasiIstirahat
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sesi belajar selesai! " + pesanIstirahat,
		"data": gin.H{
			"sesi_id":          sesi.ID,
			"mata_kuliah":      sesi.MataKuliah,
			"durasi_belajar":   sesi.DurasiBelajar,
			"sesi_keberapa":    sesi.SesiKeberapa,
			"jumlah_sesi":      sesi.JumlahSesi,
			"total_pause":      sesi.TotalPause,
			"selesai_at":       now,
			"istirahat": gin.H{
				"durasi":  durasiIstirahat,
				"pesan":   pesanIstirahat,
			},
		},
	})
}