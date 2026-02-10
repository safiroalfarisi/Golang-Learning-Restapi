package main

import (
	"log"
	"os"

	"golang-learning-restapi/config"
	"golang-learning-restapi/models"
	"golang-learning-restapi/routes"
	"golang-learning-restapi/services" // Pastikan package services diimport

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Environment Variables (.env)
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan env system")
	}

	// 2. Koneksi ke Database
	config.ConnectDB()

	// 3. Jalankan Auto Migration
	config.DB.AutoMigrate(&models.Employee{})
	log.Println("✅ Database Migration Successful")

	// --- PENAMBAHAN DI SINI ---
	// 4. Jalankan Cron Job (Challenge Task)
	// Memanggil fungsi dari services/cron.service.go
	services.InitReportJob() 
	// ---------------------------

	// 5. Inisialisasi Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Employee Management API v1.0",
	})

	// 6. Middleware Logger
	app.Use(logger.New())

	// 7. Setup Routes
	routes.SetupRoutes(app)

	// 8. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}