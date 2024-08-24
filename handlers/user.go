package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// hashPassword hashes the user's password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// verifyPassword checks if the password matches a specific pattern (e.g., minimum length)
func verifyPassword(password string) bool {
	// Example: Password must be at least 8 characters long
	if len(password) < 8 {
		return false
	}
	// Optionally add more complex rules like regex validation for special characters, etc.
	return true
}

// verifyUsername checks if the username meets specific requirements
func verifyUsername(username string) bool {
	// Example: Username must be alphanumeric and between 3 to 20 characters long
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]{3,20}$`, username)
	return matched

}

var users = map[string]string{}

func Register(w http.ResponseWriter, r *http.Request) {
	// Serve the registration page on GET request
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("frontend/register.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			log.Printf("Error loading template: %v\n", err)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	// Handle the registration on POST request
	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			log.Printf("Error parsing form data: %v\n", err)
			return
		}

		// Get username and password from the form
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Check if the username already exists
		if _, exists := users[username]; exists {
			http.Error(w, "Username already exists. Please choose a different username.", http.StatusConflict)
			return
		}

		// Verify the username
		if !verifyUsername(username) {
			http.Error(w, "Invalid username. It must be alphanumeric and 3-20 characters long.", http.StatusBadRequest)
			return
		}

		// Verify the password
		if !verifyPassword(password) {
			http.Error(w, "Invalid password. It must be at least 8 characters long.", http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := hashPassword(password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			log.Printf("Error hashing password: %v\n", err)
			return
		}

		// Store the new user's username and hashed password
		users[username] = hashedPassword

		// Simulate successful registration
		log.Printf("User %s has registered successfully", username)
		fmt.Fprintf(w, "User %s registered successfully!", username)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmp, err := template.ParseFiles("frontend/login.html")
		if err != nil {
			http.Error(w, "Error fetching template file", http.StatusInternalServerError)
			log.Printf("Error loading template: %v\n", err)
			return
		}
		tmp.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		// Fetch form values
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Retrieve the stored hashed password for the user
		storedHashedPassword, ok := users[username]
		if !ok {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Compare hashed password with the submitted password
		err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Login successful
		fmt.Fprintf(w, "Welcome, %s! You have successfully logged in", username)
		log.Printf("User %s logged in successfully", username)
	}
}
