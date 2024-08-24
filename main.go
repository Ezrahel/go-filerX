package main

import (
	"fmt"
	"net/http"

	"github.com/ezrahel/pdfupload/handlers"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	handlers.DB()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))
	http.HandleFunc("/create", handlers.CreatePDF)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/listpdf", handlers.GetPDFs)
	http.HandleFunc("/login", handlers.LoginUser)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Listening on port 8080 ")
}
