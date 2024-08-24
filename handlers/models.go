package handlers

import (
	"path/filepath"
	"time"
)

type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdat"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PdfResource struct {
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	PdfFile   string    `json:"pdffile"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

type PdfUpload struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	File        filepath.WalkFunc `json:"file"`
}
