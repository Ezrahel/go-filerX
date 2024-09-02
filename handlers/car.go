package handlers

import (
	"fmt"
	"log"
	"time"
)

func Car(make string, model string, year int) error {
	// Open the database connection
	db, err := DB()
	if err != nil {
		log.Printf("Database connection error: %v\n", err)
		return fmt.Errorf("Database connection error: %v", err)
	}

	// Ensure the database connection is closed after the program finishes using it
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Printf("Error closing database: %v\n", cerr)
		}
	}()

	// Create the insertion query
	query := "INSERT INTO cars (make, model, year, created_at) VALUES (?, ?, ?, ?)"

	// Get the current timestamp for created_at
	createdAt := time.Now()

	// Execute the query
	_, err = db.Exec(query, make, model, year, createdAt)
	if err != nil {
		log.Printf("Error inserting car into database: %v\n", err)
		return fmt.Errorf("Error inserting car into database: %v", err)
	}

	log.Printf("Car %s %s inserted successfully!", make, model)
	return nil
}
