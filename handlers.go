package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type PageData struct {
    GPA   float64
    Error string
}

const maxUploadSize = 5 * 1024 * 1024 // 5 MB

// main page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		return
	}
	data := PageData{}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

// uploadTranscriptHandler handles the file upload and GPA calculation
func uploadTranscriptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			log.Printf("Error parsing template: %v", err)
			http.Error(w, "Error loading page", http.StatusInternalServerError)
			return
		}
		data := PageData{}
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
		}
		return
	}

	if r.Method == "POST" {
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		file, _, err := r.FormFile("file")
		if err != nil {
			log.Printf("Error parsing file: %v", err)
			renderTemplateWithError(w, "Invalid file upload. Please try again.")
			return
		}
		defer file.Close()

		// Create a temporary file to store the uploaded transcript
		tempFile, err := os.CreateTemp("", "transcript-*.pdf")
		if err != nil {
			log.Printf("Could not create temporary file: %v", err)
			renderTemplateWithError(w, "An error occurred while processing your request.")
			return
		}
		defer tempFile.Close()
		defer os.Remove(tempFile.Name())

		// Write the uploaded content to the temp file
		_, err = tempFile.ReadFrom(file)
		if err != nil {
			log.Printf("Error saving uploaded file: %v", err)
			renderTemplateWithError(w, "An error occurred while saving your file.")
			return
		}

		// Calculate GPA from the PDF
		_, gpa, err := CalculateGPA(tempFile.Name())
		if err != nil {
			log.Printf("Error processing PDF: %v", err)
			renderTemplateWithError(w, "Failed to process your transcript. Ensure it's a valid Waterloo Works Transcript.")
			return
		}

		// Render the index.html template with the GPA result
		tmpl, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			log.Printf("Error parsing template: %v", err)
			http.Error(w, "Error loading page", http.StatusInternalServerError)
			return
		}
		data := PageData{
			GPA: gpa,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
		}
	}
}

// renderTemplateWithError renders the index.html template with an error message
func renderTemplateWithError(w http.ResponseWriter, errorMessage string) {
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		return
	}
	data := PageData{
		Error: errorMessage,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
