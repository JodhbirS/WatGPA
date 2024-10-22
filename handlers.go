package main

import (
    "html/template"
    "log"
    "net/http"
    "os"
)

const maxUploadSize = 5 * 1024 * 1024 // 5 MB

// main page
func homeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("./templates/index.html")
    if err != nil {
        http.Error(w, "Error loading page", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

// uploadTranscriptHandler handles the file upload
func uploadTranscriptHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        // Render the index.html template without GPA
        tmpl, err := template.ParseFiles("./templates/index.html")
        if err != nil {
            http.Error(w, "Error loading page", http.StatusInternalServerError)
            return
        }
        tmpl.Execute(w, nil)
        return
    }

    if r.Method == "POST" {
        r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

        file, _, err := r.FormFile("file")
        if err != nil {
            log.Printf("Error parsing file: %v", err)
            http.Error(w, "Invalid file upload. Please try again.", http.StatusBadRequest)
            return
        }
        defer file.Close()

        tempFile, err := os.CreateTemp("", "transcript-*.pdf")
        if err != nil {
            log.Printf("Could not create temporary file: %v", err)
            http.Error(w, "An error occurred while processing your request.", http.StatusInternalServerError)
            return
        }
        defer tempFile.Close()
        defer os.Remove(tempFile.Name()) // Clean up the temp file

        // Write the uploaded content to the temp file
        _, err = tempFile.ReadFrom(file)
        if err != nil {
            log.Printf("Error saving uploaded file: %v", err)
            http.Error(w, "An error occurred while processing your request.", http.StatusInternalServerError)
            return
        }

        // Calculate GPA from the PDF
        _, gpa, err := CalculateGPA(tempFile.Name())
        if err != nil {
            log.Printf("Error processing PDF: %v", err)
            http.Error(w, "An error occurred while processing your transcript.", http.StatusInternalServerError)
            return
        }

        // Render the index.html template with the GPA result
        tmpl, err := template.ParseFiles("./templates/index.html")
        if err != nil {
            http.Error(w, "Error loading page", http.StatusInternalServerError)
            return
        }
        data := struct {
            GPA float64
        }{
            GPA: gpa,
        }
        tmpl.Execute(w, data)
    }
}