# **WatGPA**

## **Overview**

**WatGPA** is a web application that calculates your GPA based on your unofficial transcript from the University of Waterloo. Upload your PDF transcript, and the app extracts your grades to compute your cumulative GPA using the OUAC scale.

## **Technologies Used**

- **Go (Golang)**: Backend language for processing and server handling.
- **Go Packages**:
  - [`net/http`](https://pkg.go.dev/net/http): For handling HTTP requests and responses.
  - [`github.com/gorilla/mux`](https://github.com/gorilla/mux): For routing.
  - [`github.com/dslipak/pdf`](https://github.com/dslipak/pdf): For extracting text from PDFs.
- **HTML/CSS/JavaScript**: Frontend technologies for the user interface.


### **Set Up and Installation**

1. **Clone the Repository**:

    ```bash
    git clone https://github.com/JodhbirS/WatGPA.git
    cd WatGPA
    ```

2. **Download Dependencies**:

    ```bash
    go mod download
    ```

4. **Build and Run the App**:

    ```bash
    go build -o watgpa
    ./watgpa
    ```
    
    The server will start at [http://localhost:8080](http://localhost:8080).
