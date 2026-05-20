package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nama        string     `json:"nama"`
	Email       string     `json:"email" gorm:"unique"`
	Password    string     `json:"-"`
	Universitas string     `json:"universitas"`
	Jurusan     string     `json:"jurusan"`
	Avatar      string     `json:"avatar"`
	IsVerified  bool       `json:"is_verified" gorm:"default:false"`
	VerifiedAt  *time.Time `json:"verified_at"`
}

type OTP struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	Kode      string    `json:"kode"`
	ExpiredAt time.Time `json:"expired_at"`
	Percobaan int       `json:"percobaan" gorm:"default:0"`
	IsUsed    bool      `json:"is_used" gorm:"default:false"`
}