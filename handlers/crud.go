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

func InstancePDF(title string, author string, mypdf os.File) error {
	db, err := DB()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)

	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v\n", err)
		}
	}()

	query := "INSERT INTO pdf (title, author, mypdf)"
	exe, err := db.Exec(query, title, author, mypdf)
	if err != nil {
		log.Printf("Error inserting to db: %v", err)
		return fmt.Errorf("Error inserting to db: %v", err)
	}
	log.Printf("PDF %s %s %v inserted successfully!", title, author, exe)
	return nil
}
func CreatePDF(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("frontend/homepage.html")
		if err != nil {
			http.Error(w, "could not parse template", http.StatusBadRequest)
		}
		tmpl.Execute(w, tmpl)
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(50 << 20)
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}
		title := r.FormValue("title")
		author := r.FormValue("author")
		file, handler, err := r.FormFile("pdffile")
		if err != nil {
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			log.Printf("Error retrieving the file: %v\n", err)
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
		details := User{
			Username: username,
		}
		pdf := PdfResource{
			Title:   title,
			PdfFile: filePath,
			Author:  author,
		}

		fmt.Printf("Saved PDF: %+v\n", pdf)

		w.WriteHeader(http.StatusCreated)
		fmt.Println(details)
		tmpl, err := template.ParseFiles("frontend/listpdf.html")
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
	// var getRes []PdfResource
	// for _, v:= range getRes{

	// }
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
