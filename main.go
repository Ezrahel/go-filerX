package main

import (
	"log"
	"net/http"

	"github.com/ezrahel/pdfupload/handlers"
)

func main() {
	if err := handlers.InitStorage(); err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/create", handlers.CreatePDF)
	http.HandleFunc("/listpdf", handlers.GetPDFs)
	http.HandleFunc("/delete", handlers.DeletePDF)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.LoginUser)

	log.Println("Listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
