package database

import (
	"log"
	"sync"
	"time"
)

// In-memory database implementation
var (
	users   = make(map[string]*User)
	resumes = make(map[string]*Resume)
	mutex   = &sync.RWMutex{}
)

// User represents a user in the in-memory database
type User struct {
	ID        string
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
}

// Resume represents a resume in the in-memory database
type Resume struct {
	ID         string
	UserID     string
	Filename   string
	Content    string
	Keywords   []string
	Skills     []string
	Education  []string
	Experience []string
	UploadedAt time.Time
}

// InitDB initializes the in-memory database
func InitDB() {
	log.Println("Using in-memory database (SQLite not available)")
	// Nothing to do for in-memory database
}

// CloseDB closes the database connection
func CloseDB() {
	// Nothing to do for in-memory database
}

// SaveUser saves a user to the in-memory database
func SaveUser(user *User) error {
	mutex.Lock()
	defer mutex.Unlock()
	users[user.ID] = user
	return nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(id string) (*User, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	user, exists := users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*User, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, nil
}

// SaveResume saves a resume to the in-memory database
func SaveResume(resume *Resume) error {
	mutex.Lock()
	defer mutex.Unlock()
	resumes[resume.ID] = resume
	return nil
}

// GetResumeByID retrieves a resume by ID
func GetResumeByID(id string) (*Resume, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	resume, exists := resumes[id]
	if !exists {
		return nil, nil
	}
	return resume, nil
}

// GetResumesByUserID retrieves all resumes for a user
func GetResumesByUserID(userID string) ([]*Resume, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	var userResumes []*Resume
	for _, resume := range resumes {
		if resume.UserID == userID {
			userResumes = append(userResumes, resume)
		}
	}
	return userResumes, nil
}
