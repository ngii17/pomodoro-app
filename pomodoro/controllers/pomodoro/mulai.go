package pomodoro

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Mulai(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input struct {
		MataKuliah             string `json:"mata_kuliah" binding:"required"`
		Metode                 string `json:"metode" binding:"required"`
		DurasiBelajar          int    `json:"durasi_belajar"`
		DurasiIstirahat        int    `json:"durasi_istirahat"`
		DurasiIstirahatPanjang int    `json:"durasi_istirahat_panjang"`
		JumlahSesi             int    `json:"jumlah_sesi"`
	}

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Validasi metode
	if input.Metode != "original" && input.Metode != "custom" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Metode harus original atau custom!",
		})
		return
	}

	// Cek apakah ada sesi yang masih aktif
	var aktivSesi models.Pomodoro
	if err := config.DB.Where("user_id = ? AND status IN ?", userID, []string{"running", "paused"}).
		First(&aktivSesi).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kamu masih punya sesi yang aktif! Selesaikan dulu sebelum mulai sesi baru.",
			"data": gin.H{
				"sesi_id": aktivSesi.ID,
				"status":  aktivSesi.Status,
			},
		})
		return
	}

	// Set durasi berdasarkan metode
	durasiBelajar := input.DurasiBelajar
	durasiIstirahat := input.DurasiIstirahat
	durasiIstirahatPanjang := input.DurasiIstirahatPanjang
	jumlahSesi := input.JumlahSesi

	if input.Metode == "original" {
		durasiBelajar = 25
		durasiIstirahat = 5
		durasiIstirahatPanjang = 15
		jumlahSesi = 4
	} else {
		// Validasi input custom
		if durasiBelajar <= 0 || durasiIstirahat <= 0 || durasiIstirahatPanjang <= 0 || jumlahSesi <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Semua durasi dan jumlah sesi harus diisi untuk metode custom!",
			})
			return
		}
	}

	// Buat sesi baru
	now := time.Now()
	sesi := models.Pomodoro{
		UserID:                 userID.(uint),
		MataKuliah:             input.MataKuliah,
		Metode:                 input.Metode,
		DurasiBelajar:          durasiBelajar,
		DurasiIstirahat:        durasiIstirahat,
		DurasiIstirahatPanjang: durasiIstirahatPanjang,
		JumlahSesi:             jumlahSesi,
		SesiKeberapa:           1,
		Status:                 "running",
		MulaiAt:                &now,
		TotalPause:             0,
	}

	if err := config.DB.Create(&sesi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat sesi!",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Sesi belajar dimulai! Semangat! 💪",
		"data":    sesi,
	})
}