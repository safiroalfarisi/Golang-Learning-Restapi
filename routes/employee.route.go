package routes

import (
	"golang-learning-restapi/handlers"
	"golang-learning-restapi/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Auth Routes
	app.Post("/api/register", handlers.Register)
	app.Post("/api/login", handlers.Login)

	// Employee Routes (Terproteksi JWT)
	api := app.Group("/api/employees", middlewares.AuthenticateJWT)
	api.Get("/", handlers.GetAllEmployees)
	api.Post("/", handlers.CreateEmployee)
}