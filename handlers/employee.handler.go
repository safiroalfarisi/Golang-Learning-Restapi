package handlers

import (
	"golang-learning-restapi/config"
	"golang-learning-restapi/models"
	"github.com/gofiber/fiber/v2"
)

// 1. Get all employees
func GetAllEmployees(c *fiber.Ctx) error {
	var employees []models.Employee
	config.DB.Find(&employees)
	return c.JSON(employees)
}

// 2. Create an employee
func CreateEmployee(c *fiber.Ctx) error {
	employee := new(models.Employee)
	if err := c.BodyParser(employee); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	// Simpan ke database
	if err := config.DB.Create(&employee).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	// Panggil service email secara asynchronous (menggunakan goroutine)
	// agar tidak memperlambat respon API
	go services.SendWelcomeEmail(employee.Email, employee.Name)

	return c.Status(201).JSON(employee)
}