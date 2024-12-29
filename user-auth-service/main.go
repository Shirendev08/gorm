package main

import (
	"fmt"
	"net/http"
	"time"
	"user-auth-service/middleware"
	"user-auth-service/models"
	"user-auth-service/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

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

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow Nuxt frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // CORS max age
	}
	router.Use(cors.New(corsConfig)) // Set up CORS middleware

	// Register routes
	router.POST("/register", registerHandler)
	router.POST("/login", loginHandler)
	router.GET("/users", getUsersHandler)

	// Use the middleware for protected routes
	protected := router.Group("/")
	protected.Use(middleware.TokenValidationMiddleware(DB))
	{
		protected.GET("/protected", protectedHandler) // Add your protected routes here
	}

	// Start server
	router.Run(":8080")
}

func registerHandler(c *gin.Context) {
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

func getUsersHandler(c *gin.Context) {
	// Fetch all users from the database
	var users []models.User
	if err := DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Return the list of users
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func protectedHandler(c *gin.Context) {
	// This is a protected route, only accessible after token validation
	c.JSON(http.StatusOK, gin.H{
		"message": "This is a protected route",
	})
}
