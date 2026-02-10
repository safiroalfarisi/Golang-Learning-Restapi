package handlers

import (
	"golang-learning-restapi/config"
	"golang-learning-restapi/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
	"os"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	employee := models.Employee{
		Name:     data["name"],
		Email:    data["email"],
		Position: data["position"],
		Password: string(password),
		Role:     "staff", // Default role
	}

	config.DB.Create(&employee)
	return c.JSON(employee)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var employee models.Employee
	config.DB.Where("email = ?", data["email"]).First(&employee)

	if employee.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(data["password"])); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Password salah"})
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   employee.ID,
		"role": employee.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

