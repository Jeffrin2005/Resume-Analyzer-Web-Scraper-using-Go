package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"webscrapper/database"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Password is never returned in JSON
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUser creates a new user in the database
func CreateUser(username, password, email string) (*User, error) {
	// Generate a unique ID
	id := uuid.New().String()
	
	// Hash the password
	hashedPassword := hashPassword(password)
	
	// Create user object
	user := &database.User{
		ID:        id,
		Username:  username,
		Password:  hashedPassword,
		Email:     email,
		CreatedAt: time.Now(),
	}
	
	// Save to in-memory database
	err := database.SaveUser(user)
	if err != nil {
		return nil, err
	}
	
	return &User{
		ID:        id,
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
	}, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(id string) (*User, error) {
	dbUser, err := database.GetUserByID(id)
	if err != nil || dbUser == nil {
		return nil, err
	}
	
	return &User{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
	}, nil
}

// GetUserByUsername retrieves a user by their username
func GetUserByUsername(username string) (*User, error) {
	dbUser, err := database.GetUserByUsername(username)
	if err != nil || dbUser == nil {
		return nil, err
	}
	
	return &User{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		Password:  dbUser.Password,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
	}, nil
}

// VerifyPassword checks if the provided password matches the stored hash
func VerifyPassword(user *User, password string) bool {
	hashedInput := hashPassword(password)
	return user.Password == hashedInput
}

// hashPassword creates a SHA-256 hash of the password
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
