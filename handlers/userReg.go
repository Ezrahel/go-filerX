package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// hashPassword hashes the given password using bcrypt.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func Register(w http.ResponseWriter, r *http.Request) {
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

	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			log.Printf("Error parsing form data: %v\n", err)
			return
		}

		username := r.FormValue("username")
		password_hash := r.FormValue("password")
		created_at := time.Now()

		// Hash the password
		hashedPassword, err := hashPassword(password_hash)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			log.Printf("Error hashing password: %v\n", err)
			return
		}

		// Open the database connection
		db, err := DB()
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			log.Printf("Database connection error: %v\n", err)
			return
		}

		// Make sure to close the database connection after use
		defer func() {
			if cerr := db.Close(); cerr != nil {
				log.Printf("Error closing database: %v\n", cerr)
			}
		}()

		// Execute the query
		query := "INSERT INTO pdf (username, password_hash, created_at) VALUES (?, ?, ?)"
		_, err = db.Exec(query, username, hashedPassword, created_at)
		if err != nil {
			http.Error(w, "Error inserting into database", http.StatusInternalServerError)
			log.Printf("Error inserting user: %v\n", err)
			return
		}

		// Success message

		http.Redirect(w, r, "/create", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}
