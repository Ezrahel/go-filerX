package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func Home(w http.ResponseWriter, r *http.Request) {
	name := "Israel"
	log.Printf("%v", name)
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("frontend/homepage.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			log.Printf("Error loading template: %v\n", err)
			return
		}
		tmpl.Execute(w, name)
		return
	}
}

func CreatePDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	} else {
		_, err := template.ParseFiles("frontend/homepage.html")
		if err != nil {
			http.Error(w, "could not parse template", http.StatusBadRequest)
		}
	}

	// Parse the form data
	err := r.ParseMultipartForm(50 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Extract form values
	title := r.FormValue("title")
	author := r.FormValue("author")

	// Extract the file
	file, handler, err := r.FormFile("pdfFile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath := filepath.Join("uploads", handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = dst.ReadFrom(file)
	if err != nil {
		http.Error(w, "Error writing the file", http.StatusInternalServerError)
		return
	}

	// Create a new PdfResource instance
	pdf := PdfResource{
		Title:   title,
		PdfFile: filePath,
		Author:  author,
	}

	// Save the PdfResource to the database or any storage
	// For demonstration, we will just print it
	fmt.Printf("Saved PDF: %+v\n", pdf)

	// Return a response to the client
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "PDF created successfully")
	tmpl, err := template.ParseFiles("frontend/homepage.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, pdf)
	if err != nil {
		http.Error(w, "error executing template", http.StatusInternalServerError)
		return
	}
}

func DeletePDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error passing form", http.StatusInternalServerError)
		return
	}
	fileName := r.FormValue("fileName")
	if fileName == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}
	filePath := filepath.Join("uploads", fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	fmt.Printf("Deleted PDF: %s\n", fileName)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "PDF deleted successfully")
}

func GetAuthor(author string) string {
	var pdf []PdfResource
	for _, v := range pdf {
		fmt.Printf("%v", v.Author)
		author = v.Author
	}
	return author
}

func GetPDF() ([]PdfResource, error) {
	var documents []PdfResource
	err := filepath.Walk("uploads", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			document := PdfResource{
				Title:   info.Name(),
				PdfFile: path,
				Author:  GetAuthor(path),
			}
			documents = append(documents, document)
		}
		return nil
	})
	return documents, err
}

func GetPDFs(w http.ResponseWriter, r *http.Request) {
	pdfs, err := GetPDF()
	if err != nil {
		http.Error(w, "Error getting PDFs", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("frontend/pdfs.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, pdfs)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
