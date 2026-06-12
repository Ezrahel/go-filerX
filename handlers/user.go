package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func verifyPassword(password string) bool {
	return len(password) >= 8
}

func verifyUsername(username string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]{3,20}$`, username)
	return matched
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmp, err := template.ParseFiles("frontend/login.html")
		if err != nil {
			http.Error(w, "error loading template", http.StatusInternalServerError)
			log.Printf("login template load error: %v", err)
			return
		}
		_ = tmp.Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "error parsing form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	storedHashedPassword, ok := getUserHash(username)
	if !ok {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password)); err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Welcome, %s! You have successfully logged in.", username)
	log.Printf("user %s logged in successfully", username)
}
