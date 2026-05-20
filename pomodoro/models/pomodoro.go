package models

import (
	"time"

	"gorm.io/gorm"
)

type Pomodoro struct {
	gorm.Model
	UserID              uint       `json:"user_id"`
	MataKuliah          string     `json:"mata_kuliah"`
	Metode              string     `json:"metode"`              // original / custom
	DurasiBelajar       int        `json:"durasi_belajar"`      // dalam menit
	DurasiIstirahat     int        `json:"durasi_istirahat"`    // istirahat pendek
	DurasiIstirahatPanjang int     `json:"durasi_istirahat_panjang"`
	JumlahSesi          int        `json:"jumlah_sesi"`         // sesi sebelum istirahat panjang
	SesiKeberapa        int        `json:"sesi_keberapa"`       // sesi ke berapa sekarang
	Status              string     `json:"status"`              // running/paused/completed/skipped
	MulaiAt             *time.Time `json:"mulai_at"`
	SelesaiAt           *time.Time `json:"selesai_at"`
	TotalPause          int        `json:"total_pause"`         // berapa kali di-pause
	User                User       `json:"user" gorm:"foreignKey:UserID"`
}