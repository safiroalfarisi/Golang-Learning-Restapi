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
	api.Get("/", handlers.GetAllEmployees)          // GET all employees
	api.Post("/", handlers.CreateEmployee)          // POST create new employee
	api.Get("/:id", handlers.GetEmployeeByID)       // GET employee by ID
	api.Put("/:id", handlers.UpdateEmployee)        // PUT update employee
	api.Delete("/:id", handlers.DeleteEmployee)     // DELETE employee
	api.Post("/import", handlers.ImportEmployees)   // POST import from CSV
	api.Get("/export", handlers.ExportEmployees)    // GET export to CSV
}