package models

import (
	"time"

	"github.com/google/uuid"
	"webscrapper/database"
)

// Resume represents a parsed resume in the system
type Resume struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Filename   string    `json:"filename"`
	Content    string    `json:"content"`
	Keywords   []string  `json:"keywords"`
	Skills     []string  `json:"skills"`
	Education  []string  `json:"education"`
	Experience []string  `json:"experience"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// SaveResume saves a resume to the database
func SaveResume(userID, filename, content string, keywords, skills, education, experience []string) (*Resume, error) {
	// Generate a unique ID
	id := uuid.New().String()
	
	// Create resume object
	resume := &database.Resume{
		ID:         id,
		UserID:     userID,
		Filename:   filename,
		Content:    content,
		Keywords:   keywords,
		Skills:     skills,
		Education:  education,
		Experience: experience,
		UploadedAt: time.Now(),
	}
	
	// Save to in-memory database
	err := database.SaveResume(resume)
	if err != nil {
		return nil, err
	}
	
	return &Resume{
		ID:         id,
		UserID:     userID,
		Filename:   filename,
		Content:    content,
		Keywords:   keywords,
		Skills:     skills,
		Education:  education,
		Experience: experience,
		UploadedAt: time.Now(),
	}, nil
}

// GetResumeByID retrieves a resume by its ID
func GetResumeByID(id string) (*Resume, error) {
	dbResume, err := database.GetResumeByID(id)
	if err != nil || dbResume == nil {
		return nil, err
	}
	
	return &Resume{
		ID:         dbResume.ID,
		UserID:     dbResume.UserID,
		Filename:   dbResume.Filename,
		Content:    dbResume.Content,
		Keywords:   dbResume.Keywords,
		Skills:     dbResume.Skills,
		Education:  dbResume.Education,
		Experience: dbResume.Experience,
		UploadedAt: dbResume.UploadedAt,
	}, nil
}

// GetResumesByUserID retrieves all resumes for a specific user
func GetResumesByUserID(userID string) ([]*Resume, error) {
	dbResumes, err := database.GetResumesByUserID(userID)
	if err != nil {
		return nil, err
	}
	
	var resumes []*Resume
	for _, dbResume := range dbResumes {
		resumes = append(resumes, &Resume{
			ID:         dbResume.ID,
			UserID:     dbResume.UserID,
			Filename:   dbResume.Filename,
			Content:    dbResume.Content,
			Keywords:   dbResume.Keywords,
			Skills:     dbResume.Skills,
			Education:  dbResume.Education,
			Experience: dbResume.Experience,
			UploadedAt: dbResume.UploadedAt,
		})
	}
	
	return resumes, nil
}
