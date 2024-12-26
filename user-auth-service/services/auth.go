package services

import (
	"errors"
	"time"
	"user-auth-service/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtKey = []byte("your_secret_key")

// Check if a user or email already exists
func CheckUserExists(db *gorm.DB, username, email string) bool {
	var user models.User
	if err := db.Where("username = ?", username).Or("email = ?", email).First(&user).Error; err == nil {
		return true
	}
	return false
}

// RegisterUser creates a new user
func RegisterUser(db *gorm.DB, username, email, password string) (*models.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create the user
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// LoginUser authenticates a user and generates a JWT token
func LoginUser(db *gorm.DB, email, password string) (*models.User, string, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, "", errors.New("user not found")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid password")
	}

	// Generate JWT token
	token, err := generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

// Generate JWT token
func generateToken(user models.User) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   user.Email,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
