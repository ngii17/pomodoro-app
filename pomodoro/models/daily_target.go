package models

import (
	"time"

	"gorm.io/gorm"
)

type DailyTarget struct {
	gorm.Model
	UserID             uint          `json:"user_id"`
	Tanggal            time.Time     `json:"tanggal"`
	Status             string        `json:"status" gorm:"default:not_started"`
	MulaiAt            *time.Time    `json:"mulai_at"`
	SelesaiAt          *time.Time    `json:"selesai_at"`
	TotalPauseDuration int           `json:"total_pause_duration"` // dalam menit
	PausedAt           *time.Time    `json:"paused_at"`
	User               User          `json:"user" gorm:"foreignKey:UserID"`
	Matkul             []TargetMatkul `json:"matkul" gorm:"foreignKey:TargetID"`
}

type TargetMatkul struct {
	gorm.Model
	TargetID        uint   `json:"target_id"`
	MataKuliah      string `json:"mata_kuliah"`
	JumlahSesi      int    `json:"jumlah_sesi"`
	SesiSelesai     int    `json:"sesi_selesai" gorm:"default:0"`
	DurasiIstirahat int    `json:"durasi_istirahat"`
	Status          string `json:"status" gorm:"default:not_started"`
	Urutan          int    `json:"urutan"`
}