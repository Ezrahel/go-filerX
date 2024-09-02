package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// func hashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	return string(bytes), err
// }

func verifyPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

func verifyUsername(username string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]{3,20}$`, username)
	return matched

}

var users = map[string]string{}

// func Register(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		tmpl, err := template.ParseFiles("frontend/register.html")
// 		if err != nil {
// 			http.Error(w, "Error loading template", http.StatusInternalServerError)
// 			log.Printf("Error loading template: %v\n", err)
// 			return
// 		}
// 		tmpl.Execute(w, nil)
// 		return
// 	}

// 	if r.Method == http.MethodPost {
// 		err := r.ParseForm()
// 		if err != nil {
// 			http.Error(w, "Error parsing form data", http.StatusBadRequest)
// 			log.Printf("Error parsing form data: %v\n", err)
// 			return
// 		}

// 		username := r.FormValue("username")
// 		password := r.FormValue("password")

// 		if _, exists := users[username]; exists {
// 			http.Error(w, "Username already exists. Please choose a different username.", http.StatusConflict)
// 			return
// 		}

// 		if !verifyUsername(username) {
// 			http.Error(w, "Invalid username. It must be alphanumeric and 3-20 characters long.", http.StatusBadRequest)
// 			return
// 		}

// 		if !verifyPassword(password) {
// 			http.Error(w, "Invalid password. It must be at least 8 characters long.", http.StatusBadRequest)
// 			return
// 		}

// 		hashedPassword, err := hashPassword(password)
// 		if err != nil {
// 			http.Error(w, "Error hashing password", http.StatusInternalServerError)
// 			log.Printf("Error hashing password: %v\n", err)
// 			return
// 		}

// 		users[username] = hashedPassword

// 		log.Printf("User %s has registered successfully", username)
// 		//fmt.Fprintf(w, "User %s registered successfully!", username)
// 		tmp, err := template.ParseFiles("frontend/listpdf.html")
// 		if err != nil {
// 			log.Fatalf("Error %v while parsing template file", err)
// 		}
// 		tmp.Execute(w, tmp)
// 	} else {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 	}
// }

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmp, err := template.ParseFiles("frontend/login.html")
		if err != nil {
			http.Error(w, "Error fetching template file", http.StatusInternalServerError)
			log.Printf("Error loading template: %v\n", err)
			return
		}
		tmp.Execute(w, tmp)
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		storedHashedPassword, ok := users[username]
		if !ok {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, "Welcome, %s! You have successfully logged in", username)
		log.Printf("User %s logged in successfully", username)
	}
}
