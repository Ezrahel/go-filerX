package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ezrahel/pdfupload/handlers"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	handlers.DB()
	err := handlers.Car("Lexus", "GLK", 2024)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		log.Println("Car successfully inserted!")
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))
	http.HandleFunc("/create", handlers.CreatePDF)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/listpdf", handlers.GetPDFs)
	http.HandleFunc("/login", handlers.LoginUser)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Listening on port 8080 ")
}
