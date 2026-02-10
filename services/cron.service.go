package services

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"golang-learning-restapi/config"
	"golang-learning-restapi/models"
)

func InitReportJob() {
	// v3 menggunakan scheduler dengan dukungan detik (opsional) atau standar 5 digit
	c := cron.New()

	// Challenge Task: Contoh membuat laporan/log jumlah karyawan setiap menit
	// Format: "menit jam hari bulan hari-dalam-minggu"
	_, err := c.AddFunc("* * * * *", func() {
		var count int64
		config.DB.Model(&models.Employee{}).Count(&count)
		
		fmt.Printf(" [CRON] Log: Total karyawan saat ini adalah %d\n", count)
	})

	if err != nil {
		fmt.Println("❌ Gagal menjadwalkan Cron Job:", err)
		return
	}

	// Memulai scheduler di background (Goroutine)
	c.Start()
	fmt.Println("✅ Scheduled tasks (Cron Jobs) initialized.")
}