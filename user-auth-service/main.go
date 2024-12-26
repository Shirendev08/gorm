package main

import (
	"fmt"
	"net/http"
	"time"
	"user-auth-service/models"
	"user-auth-service/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// "github.com/itsjamie/gin-cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// func enableCors(w *http.ResponseWriter) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// }

func main() {
	// Initialize database connection
	dsn := "host=localhost user=postgres password=123 dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Migrate the User model
	DB.AutoMigrate(&models.User{})

	// Initialize Gin router
	router := gin.Default()
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow Nuxt frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // CORS max age
	}
	router.Use(cors.New(corsConfig))
	// Set up CORS middleware

	// Register and login routes
	router.GET("/asd/", func(c *gin.Context) {
		c.String(200, "hello cors")
	})
	router.POST("/register", registerHandler)
	router.POST("/login", loginHandler)

	// Start server
	router.Run(":8080")
}

func registerHandler(c *gin.Context) {
	// enableCors(&w)
	// Debugging line
	fmt.Println("Received registration request")

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if username or email already exists
	if services.CheckUserExists(DB, req.Username, req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User or email already exists"})
		return
	}

	// Register the user
	user, err := services.RegisterUser(DB, req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func loginHandler(c *gin.Context) {
	// Debugging line
	// w :=
	// 	enableCors(&w)
	fmt.Println("Received login request")

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Authenticate the user
	user, token, err := services.LoginUser(DB, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"username": user.Username,
			"email":    user.Email,
		},
		"token": token,
	})
}
