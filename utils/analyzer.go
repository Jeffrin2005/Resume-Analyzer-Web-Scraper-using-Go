package utils

import (
	"regexp"
	"strings"

	"github.com/jdkato/prose/v2"
)

// Common skills that might be found in resumes
var commonSkills = []string{
	"javascript", "python", "java", "c++", "c#", "go", "golang", "ruby", "php", "swift",
	"html", "css", "react", "angular", "vue", "node", "express", "django", "flask", "spring",
	"aws", "azure", "gcp", "docker", "kubernetes", "terraform", "jenkins", "git", "github",
	"sql", "mysql", "postgresql", "mongodb", "redis", "elasticsearch", "graphql", "rest api",
	"machine learning", "artificial intelligence", "data science", "big data", "data analysis",
	"tensorflow", "pytorch", "pandas", "numpy", "scikit-learn", "nlp", "computer vision",
	"agile", "scrum", "kanban", "jira", "confluence", "leadership", "teamwork", "communication",
}

// Common education keywords
var educationKeywords = []string{
	"bachelor", "master", "phd", "doctorate", "degree", "university", "college", "institute",
	"b.tech", "m.tech", "b.e.", "m.e.", "b.sc", "m.sc", "b.a.", "m.a.", "mba", "certification",
}

// AnalyzeResume extracts key information from resume text
func AnalyzeResume(text string) (keywords []string, skills []string, education []string, experience []string) {
	// Clean the text
	text = cleanText(text)
	
	// Extract keywords using NLP
	keywords = extractKeywords(text)
	
	// Extract skills
	skills = extractSkills(text)
	
	// Extract education
	education = extractEducation(text)
	
	// Extract experience
	experience = extractExperience(text)
	
	return
}

// cleanText removes extra whitespace and normalizes text
func cleanText(text string) string {
	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")
	
	// Convert to lowercase for easier matching
	text = strings.ToLower(text)
	
	return strings.TrimSpace(text)
}

// extractKeywords extracts important keywords using NLP
func extractKeywords(text string) []string {
	doc, err := prose.NewDocument(text)
	if err != nil {
		return []string{}
	}
	
	// Extract entities and nouns as keywords
	keywordMap := make(map[string]bool)
	
	// Add named entities
	for _, ent := range doc.Entities() {
		if len(ent.Text) > 3 { // Ignore very short entities
			keywordMap[strings.ToLower(ent.Text)] = true
		}
	}
	
	// Add important nouns
	for _, tok := range doc.Tokens() {
		if strings.HasPrefix(tok.Tag, "NN") && len(tok.Text) > 3 {
			keywordMap[strings.ToLower(tok.Text)] = true
		}
	}
	
	// Convert map to slice
	var keywords []string
	for k := range keywordMap {
		keywords = append(keywords, k)
	}
	
	// Limit to top 50 keywords
	if len(keywords) > 50 {
		keywords = keywords[:50]
	}
	
	return keywords
}

// extractSkills identifies skills mentioned in the resume
func extractSkills(text string) []string {
	var foundSkills []string
	
	for _, skill := range commonSkills {
		if strings.Contains(text, skill) {
			foundSkills = append(foundSkills, skill)
		}
	}
	
	return foundSkills
}

// extractEducation identifies education information
func extractEducation(text string) []string {
	var educationInfo []string
	
	// Split text into paragraphs
	paragraphs := strings.Split(text, "\n")
	
	for _, para := range paragraphs {
		// Check if paragraph contains education keywords
		isEducation := false
		for _, keyword := range educationKeywords {
			if strings.Contains(para, keyword) {
				isEducation = true
				break
			}
		}
		
		if isEducation && len(para) > 10 {
			educationInfo = append(educationInfo, strings.TrimSpace(para))
		}
	}
	
	return educationInfo
}

// extractExperience identifies work experience information
func extractExperience(text string) []string {
	var experienceInfo []string
	
	// Regular expressions for common experience patterns
	experiencePatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(experience|work|employment|job).*?:`),
		regexp.MustCompile(`(?i)(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[\s\.,]+\d{4}\s*[-–—]\s*(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec|present)`),
		regexp.MustCompile(`(?i)\d{4}\s*[-–—]\s*(\d{4}|present)`),
	}
	
	// Split text into paragraphs
	paragraphs := strings.Split(text, "\n")
	
	for _, para := range paragraphs {
		// Check if paragraph contains experience patterns
		isExperience := false
		for _, pattern := range experiencePatterns {
			if pattern.MatchString(para) {
				isExperience = true
				break
			}
		}
		
		if isExperience && len(para) > 10 {
			experienceInfo = append(experienceInfo, strings.TrimSpace(para))
		}
	}
	
	return experienceInfo
}
