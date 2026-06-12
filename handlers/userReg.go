package handlers

import (
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("frontend/register.html")
		if err != nil {
			http.Error(w, "error loading template", http.StatusInternalServerError)
			log.Printf("register template load error: %v", err)
			return
		}
		_ = tmpl.Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if !verifyUsername(username) || !verifyPassword(password) {
		http.Error(w, "invalid username or password", http.StatusBadRequest)
		return
	}

	if usernameTaken(username) {
		http.Error(w, "username already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, "unable to hash password", http.StatusInternalServerError)
		return
	}

	storeUser(username, hashedPassword)
	if err := saveUsersSnapshot("users.json"); err != nil {
		log.Printf("warning: unable to persist users snapshot: %v", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
