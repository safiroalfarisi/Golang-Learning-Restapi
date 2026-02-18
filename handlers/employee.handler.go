package handlers

import (
	"encoding/csv"
	"fmt"
	"strings"

	"golang-learning-restapi/config"
	"golang-learning-restapi/models"
	"golang-learning-restapi/services"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(employee.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to hash password"})
	}
	employee.Password = string(hashedPassword)

	// Panggil service email secara asynchronous (menggunakan goroutine)
	// agar tidak memperlambat respon API
	go services.SendWelcomeEmail(employee.Email, employee.Name)

	return c.Status(201).JSON(employee)
}

// 3. Get employee by ID
func GetEmployeeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var employee models.Employee

	if err := config.DB.First(&employee, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Employee not found"})
	}

	return c.JSON(employee)
}

// 4. Update employee
func UpdateEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	var employee models.Employee

	// Cari employee berdasarkan ID
	if err := config.DB.First(&employee, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Employee not found"})
	}

	// Parse request body
	var updateData models.Employee
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	// Update fields yang tidak kosong
	if updateData.Name != "" {
		employee.Name = updateData.Name
	}
	if updateData.Email != "" {
		employee.Email = updateData.Email
	}
	if updateData.Position != "" {
		employee.Position = updateData.Position
	}
	if updateData.Role != "" {
		employee.Role = updateData.Role
	}
	if updateData.Password != "" {
		// Hash password baru
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Failed to hash password"})
		}
		employee.Password = string(hashedPassword)
	}

	// Simpan perubahan
	if err := config.DB.Save(&employee).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(employee)
}

// 5. Delete employee
func DeleteEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	var employee models.Employee

	// Cari employee berdasarkan ID
	if err := config.DB.First(&employee, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Employee not found"})
	}

	// Soft delete
	if err := config.DB.Delete(&employee).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Employee deleted successfully"})
}

// 6. Import employees from CSV
func ImportEmployees(c *fiber.Ctx) error {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "No file uploaded"})
	}

	// Check file type
	if !strings.HasSuffix(file.Filename, ".csv") {
		return c.Status(400).JSON(fiber.Map{"message": "Only CSV files are allowed"})
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Cannot open file"})
	}
	defer src.Close()

	// Read CSV
	reader := csv.NewReader(src)
	records, err := reader.ReadAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Cannot read CSV file"})
	}

	var employees []models.Employee
	var errors []string

	// Skip header row and process data
	for i, record := range records[1:] {
		if len(record) < 4 {
			errors = append(errors, fmt.Sprintf("Row %d: insufficient columns", i+2))
			continue
		}

		// Hash default password
		defaultPassword := "password123"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: failed to hash password", i+2))
			continue
		}

		employee := models.Employee{
			Name:     record[0],
			Email:    record[1],
			Position: record[2],
			Role:     record[3],
			Password: string(hashedPassword),
		}

		employees = append(employees, employee)
	}

	// Bulk insert
	if len(employees) > 0 {
		if err := config.DB.Create(&employees).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed to import employees",
				"error":   err.Error(),
				"errors":  errors,
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Successfully imported %d employees", len(employees)),
		"errors":  errors,
		"count":   len(employees),
	})
}

// 7. Export employees to CSV
func ExportEmployees(c *fiber.Ctx) error {
	var employees []models.Employee
	config.DB.Find(&employees)

	// Set response headers for CSV download
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=employees.csv")

	// Create CSV content
	var csvContent strings.Builder
	csvContent.WriteString("Name,Email,Position,Role,Created At\n")

	for _, emp := range employees {
		csvContent.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s\n",
			emp.Name,
			emp.Email,
			emp.Position,
			emp.Role,
			emp.CreatedAt.Format("2006-01-02 15:04:05"),
		))
	}

	return c.SendString(csvContent.String())
}