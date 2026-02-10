package services

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendWelcomeEmail(toEmail string, userName string) error {
	// Ambil konfigurasi dari .env
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Welcome to the Company!")
	
	// Isi email (Body)
	body := fmt.Sprintf("Hello %s,\n\nWelcome to our company! Your employee account has been successfully created.", userName)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// Kirim Email
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("❌ Gagal mengirim email:", err)
		return err
	}

	fmt.Printf("📧 Email sambutan terkirim ke: %s\n", toEmail)
	return nil
}