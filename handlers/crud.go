package handlers

import (
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const maxUploadSize = 50 << 20

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("frontend/homepage.html")
	if err != nil {
		http.Error(w, "error loading template", http.StatusInternalServerError)
		log.Printf("home template error: %v", err)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
	}
}

func CreatePDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "file too large or invalid form data", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("pdffile")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName, err := safeUploadName(fileHeader.Filename)
	if err != nil {
		http.Error(w, "invalid file name", http.StatusBadRequest)
		return
	}

	ext := strings.ToLower(filepath.Ext(fileName))
	allowed := map[string]bool{".pdf": true, ".doc": true, ".docx": true, ".ppt": true, ".pptx": true, ".xls": true, ".xlsx": true, ".txt": true, ".odt": true}
	if !allowed[ext] {
		http.Error(w, "unsupported file type", http.StatusBadRequest)
		return
	}

	dstPath := filepath.Join(uploadDir, fileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "failed to write file", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/listpdf", http.StatusSeeOther)
}

func DeletePDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fileName := r.FormValue("fileName")
	if fileName == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}
	path := filepath.Join(uploadDir, filepath.Base(fileName))
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		http.Error(w, "could not delete file", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/listpdf", http.StatusSeeOther)
}

func GetPDFs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	pdfs, err := listUploadedFiles()
	if err != nil {
		http.Error(w, "error reading uploads", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("frontend/listpdf.html")
	if err != nil {
		http.Error(w, "error parsing template", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, pdfs); err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
	}
}
