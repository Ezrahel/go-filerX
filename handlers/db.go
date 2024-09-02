package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "Taepryung024,"
	hostname = "127.0.0.1:3306"
	dbname   = "pdfdrive"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func DB() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, fmt.Errorf("Error %v opening database", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	// Ping the database to ensure it's connected
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Error %s pinging DB", err)
		return nil, fmt.Errorf("Error %v pinging database", err)
	}

	log.Printf("Connected to DB %s successfully\n", dbname)
	return db, nil
}
