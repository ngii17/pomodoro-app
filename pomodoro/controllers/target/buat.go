package target

import (
	"net/http"
	"time"

	"studytrack/config"
	"studytrack/models"

	"github.com/gin-gonic/gin"
)

func Buat(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input struct {
		Tanggal string `json:"tanggal" binding:"required"` // format: 2026-05-17
		Matkul  []struct {
			MataKuliah      string `json:"mata_kuliah" binding:"required"`
			JumlahSesi      int    `json:"jumlah_sesi" binding:"required"`
			DurasiIstirahat int    `json:"durasi_istirahat" binding:"required"`
		} `json:"matkul" binding:"required"`
	}

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Validasi minimal 1 mata kuliah
	if len(input.Matkul) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Minimal harus ada 1 mata kuliah!",
		})
		return
	}

	// Parse tanggal
	tanggal, err := time.Parse("2006-01-02", input.Tanggal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format tanggal tidak valid! Gunakan format: 2026-05-17",
		})
		return
	}

	// Cek apakah sudah ada target di tanggal yang sama
	var existingTarget models.DailyTarget
	if err := config.DB.Where("user_id = ? AND DATE(tanggal) = ?", userID, input.Tanggal).
		First(&existingTarget).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kamu sudah punya target di tanggal ini!",
			"data": gin.H{
				"target_id": existingTarget.ID,
			},
		})
		return
	}

	// Buat target baru
	target := models.DailyTarget{
		UserID:  userID.(uint),
		Tanggal: tanggal,
		Status:  "not_started",
	}

	if err := config.DB.Create(&target).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat target!",
		})
		return
	}

	// Simpan semua mata kuliah
	for i, m := range input.Matkul {
		matkul := models.TargetMatkul{
			TargetID:        target.ID,
			MataKuliah:      m.MataKuliah,
			JumlahSesi:      m.JumlahSesi,
			DurasiIstirahat: m.DurasiIstirahat,
			Status:          "not_started",
			Urutan:          i + 1,
		}
		config.DB.Create(&matkul)
	}

	// Ambil target lengkap dengan matkulnya
	config.DB.Preload("Matkul").First(&target, target.ID)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Target harian berhasil dibuat!",
		"data":    target,
	})
}