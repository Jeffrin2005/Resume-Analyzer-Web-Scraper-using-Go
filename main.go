package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ledongthuc/pdf"
	
	"webscrapper/database"
	"webscrapper/handlers"
	"webscrapper/utils"
)

// legacyFormHandler is the original form handler kept for backward compatibility
func legacyFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")

	file, header, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save uploaded file temporarily
	tmpFile, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		http.Error(w, "Unable to create temp file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		http.Error(w, "Unable to save uploaded file", http.StatusInternalServerError)
		return
	}

	// Extract text from PDF
	tmpFile.Seek(0, 0)
	reader, err := pdf.NewReader(tmpFile, header.Size)
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

	// Analyze the resume text
	_, skills, education, experience := utils.AnalyzeResume(text)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h2>Resume Summary for %s</h2>", name)
	fmt.Fprintf(w, "<h3>Extracted Text:</h3>")
	fmt.Fprintf(w, "<pre style='white-space: pre-wrap;'>%s</pre>", text)
	
	fmt.Fprintf(w, "<h3>Skills:</h3>")
	fmt.Fprintf(w, "<ul>")
	for _, skill := range skills {
		fmt.Fprintf(w, "<li>%s</li>", skill)
	}
	fmt.Fprintf(w, "</ul>")
	
	fmt.Fprintf(w, "<h3>Education:</h3>")
	fmt.Fprintf(w, "<ul>")
	for _, edu := range education {
		fmt.Fprintf(w, "<li>%s</li>", edu)
	}
	fmt.Fprintf(w, "</ul>")
	
	fmt.Fprintf(w, "<h3>Experience:</h3>")
	fmt.Fprintf(w, "<ul>")
	for _, exp := range experience {
		fmt.Fprintf(w, "<li>%s</li>", exp)
	}
	fmt.Fprintf(w, "</ul>")
	
	fmt.Fprint(w, `<p><a href='/form.html'>Back to Upload</a></p>`)
	fmt.Fprint(w, `<p><a href='/app'>Try our new enhanced app!</a></p>`)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, "Hello, World!")
}

// setupRoutes configures all the routes for our application
func setupRoutes() http.Handler {
	// Initialize Chi router
	r := chi.NewRouter()
	
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	
	// Public routes
	r.Group(func(r chi.Router) {
		// Serve static files
		fs := http.FileServer(http.Dir("./static"))
		r.Handle("/*", http.StripPrefix("/", fs))
		
		// Legacy routes for backward compatibility
		r.Get("/hello", helloHandler)
		r.Post("/form", legacyFormHandler)
		
		// Auth routes
		r.Post("/api/login", handlers.LoginHandler)
		r.Post("/api/register", handlers.RegisterHandler)
	})
	
	// Protected routes (require authentication)
	r.Group(func(r chi.Router) {
		// Apply auth middleware
		r.Use(utils.AuthMiddleware)
		
		// Resume routes
		r.Post("/api/resumes/upload", handlers.UploadResumeHandler)
		r.Get("/api/resumes", handlers.GetResumesHandler)
		r.Get("/api/resumes/{id}", handlers.GetResumeHandler)
	})
	
	// Serve the SPA for any routes not matched
	r.Get("/app*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("./static", "app.html"))
	})
	
	return r
}

func main() {
	// Initialize database
	database.InitDB()
	defer database.CloseDB()
	
	// Create uploads directory if it doesn't exist
	uploadsDir := "./uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		os.Mkdir(uploadsDir, 0755)
	}
	
	// Setup routes
	router := setupRoutes()
	
	// Start server
	port := "8081"
	fmt.Printf("Starting enhanced resume analyzer server at port %s\n", port)
	fmt.Printf("Access the legacy app at http://localhost:%s/\n", port)
	fmt.Printf("Access the new enhanced app at http://localhost:%s/app\n", port)
	
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server error:", err)
	}
}
