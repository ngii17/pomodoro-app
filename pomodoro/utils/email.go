package utils

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// Generate OTP 6 digit
func GenerateOTP() string {
	otp := rand.Intn(900000) + 100000
	return strconv.Itoa(otp)
}

// Kirim email OTP
func SendOTPEmail(email, nama, otp string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	port, _ := strconv.Atoi(smtpPort)

	m := gomail.NewMessage()
	m.SetHeader("From", smtpEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Kode OTP Pomodoro App")
	m.SetBody("text/html", fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 500px; margin: 0 auto;">
			<h2 style="color: #e74c3c;">🍅 Pomodoro App</h2>
			<p>Halo <strong>%s</strong>!</p>
			<p>Kode OTP kamu untuk verifikasi akun:</p>
			<div style="background: #f4f4f4; padding: 20px; text-align: center; border-radius: 8px;">
				<h1 style="color: #e74c3c; letter-spacing: 8px;">%s</h1>
			</div>
			<p style="color: #888; font-size: 12px;">Kode berlaku selama <strong>5 menit</strong>.</p>
			<p style="color: #888; font-size: 12px;">Jangan bagikan kode ini ke siapapun!</p>
		</div>
	`, nama, otp))

	d := gomail.NewDialer(smtpHost, port, smtpEmail, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil

}

// Kirim email OTP Forgot Password
func SendForgotPasswordEmail(email, nama, otp string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	port, _ := strconv.Atoi(smtpPort)

	m := gomail.NewMessage()
	m.SetHeader("From", smtpEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Reset Password Pomodoro App")
	m.SetBody("text/html", fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 500px; margin: 0 auto;">
			<h2 style="color: #e74c3c;">🍅 Pomodoro App</h2>
			<p>Halo <strong>%s</strong>!</p>
			<p>Kamu baru saja meminta reset password. Gunakan kode OTP berikut:</p>
			<div style="background: #f4f4f4; padding: 20px; text-align: center; border-radius: 8px;">
				<h1 style="color: #e74c3c; letter-spacing: 8px;">%s</h1>
			</div>
			<p style="color: #888; font-size: 12px;">Kode berlaku selama <strong>5 menit</strong>.</p>
			<p style="color: #888; font-size: 12px;">Jika kamu tidak merasa meminta reset password, abaikan email ini!</p>
		</div>
	`, nama, otp))

	d := gomail.NewDialer(smtpHost, port, smtpEmail, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}