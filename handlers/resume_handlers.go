package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/ledongthuc/pdf"
	"webscrapper/models"
	"webscrapper/utils"
)

// ResumeUploadResponse represents the response after uploading a resume
type ResumeUploadResponse struct {
	ResumeID   string   `json:"resume_id"`
	Keywords   []string `json:"keywords"`
	Skills     []string `json:"skills"`
	Education  []string `json:"education"`
	Experience []string `json:"experience"`
}

// UploadResumeHandler handles resume uploads and analysis
func UploadResumeHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from request header (set by auth middleware)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check file extension
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".pdf") {
		http.Error(w, "Only PDF files are allowed", http.StatusBadRequest)
		return
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "./uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		os.Mkdir(uploadsDir, 0755)
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%s", uuid.New().String(), header.Filename)
	filepath := filepath.Join(uploadsDir, filename)

	// Save file to disk
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Reopen file for PDF processing
	dst.Close()
	f, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "Failed to process PDF", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Extract text from PDF
	reader, err := pdf.NewReader(f, header.Size)
	if err != nil {
		http.Error(w, "Failed to read PDF file", http.StatusInternalServerError)
		return
	}

	var text string
	for pageNum := 1; pageNum <= reader.NumPage(); pageNum++ {
		page := reader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}
		content, _ := page.GetPlainText(nil)
		text += content
	}

	// Analyze resume text
	keywords, skills, education, experience := utils.AnalyzeResume(text)

	// Save resume to database
	resume, err := models.SaveResume(userID, header.Filename, text, keywords, skills, education, experience)
	if err != nil {
		http.Error(w, "Failed to save resume data", http.StatusInternalServerError)
		return
	}

	// Create response
	response := ResumeUploadResponse{
		ResumeID:   resume.ID,
		Keywords:   keywords,
		Skills:     skills,
		Education:  education,
		Experience: experience,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetResumesHandler retrieves all resumes for a user
func GetResumesHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from request header (set by auth middleware)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get resumes for user
	resumes, err := models.GetResumesByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve resumes", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resumes)
}

// GetResumeHandler retrieves a specific resume
func GetResumeHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from request header (set by auth middleware)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get resume ID from URL path
	resumeID := strings.TrimPrefix(r.URL.Path, "/api/resumes/")
	if resumeID == "" {
		http.Error(w, "Resume ID is required", http.StatusBadRequest)
		return
	}

	// Get resume
	resume, err := models.GetResumeByID(resumeID)
	if err != nil {
		http.Error(w, "Resume not found", http.StatusNotFound)
		return
	}

	// Check if resume belongs to user
	if resume.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resume)
}
