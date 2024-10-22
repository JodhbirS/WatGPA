package main

import (
    "bytes"
    "fmt"
    "regexp"
	"strings"

    "github.com/dslipak/pdf"
)

// Course represents a course with its name, grade, and GPA
type Course struct {
    Name  string
    Grade int
    GPA   float64
}

// CalculateGPA extracts courses and grades from the PDF and calculates the GPA
func CalculateGPA(pdfPath string) ([]Course, float64, error) {
    content, err := readPdf(pdfPath)
    if err != nil {
        return nil, 0.0, err
    }

    courses := extractCourses(content)

    if len(courses) == 0 {
        return nil, 0.0, fmt.Errorf("no grades found in the transcript")
    }

    var sum float64
    for i, course := range courses {
        courses[i].GPA = gradeToGpa(course.Grade)
        sum += courses[i].GPA
    }
    gpa := sum / float64(len(courses))

    return courses, gpa, nil
}

// readPdf reads the content of a PDF file and returns it as a string.
func readPdf(path string) (string, error) {
    r, err := pdf.Open(path)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    b, err := r.GetPlainText()
    if err != nil {
        return "", err
    }
    buf.ReadFrom(b)
    content := buf.String()

    // print content of the PDF for debugging
    //fmt.Println("Extracted PDF Content:", content)

    return content, nil
}

// extractCourses extracts courses and their grades from the provided text.
func extractCourses(content string) []Course {
    var courses []Course
    re := regexp.MustCompile(`(?s)([A-Z]+\s*\d{3}[A-Z]?).*?0\.50\s*(0\.\d{2})\s*([0-9]{1,3})(.*?)`)

    matches := re.FindAllStringSubmatch(content, -1)
    for _, match := range matches {
        if len(match) > 4 {
            courseName := match[1]
            earnedCreditsStr := match[2]
            gradeStr := match[3]
            postGradeText := match[4]

            var grade int
            fmt.Sscanf(gradeStr, "%d", &grade)

            var earnedCredits float64
            fmt.Sscanf(earnedCreditsStr, "%f", &earnedCredits)

            includeInGPA := false

            if earnedCredits > 0.0 {
                includeInGPA = true
            } else if strings.Contains(strings.ToLower(postGradeText), "included in average") {
                includeInGPA = true
            } else if strings.Contains(strings.ToLower(postGradeText), "not in average") {
                includeInGPA = false
            }

            if includeInGPA {
                course := Course{
                    Name:  courseName,
                    Grade: grade,
                    GPA:   gradeToGpa(grade),
                }
                courses = append(courses, course)

                //fmt.Println(course.Grade)
            }
        }
    }

    return courses
}

// GPA scale given by OUAC
func gradeToGpa(grade int) float64 {
    switch {
    case grade >= 90:
        return 4.0
    case grade >= 85:
        return 3.9
    case grade >= 80:
        return 3.7
    case grade >= 77:
        return 3.3
    case grade >= 73:
        return 3.0
    case grade >= 70:
        return 2.7
    case grade >= 67:
        return 2.3
    case grade >= 63:
        return 2.0
    case grade >= 60:
        return 1.7
    case grade >= 57:
        return 1.3
    case grade >= 53:
        return 1.0
    case grade >= 50:
        return 0.7
    default:
        return 0.0
    }
}