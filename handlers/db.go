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
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		return nil, fmt.Errorf("Error %s when opening DB\n", err)

	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return nil, fmt.Errorf(" Error %v creating database: \n", err)
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return nil, fmt.Errorf(" Error %v fetching database row: \n", err)
	}
	log.Printf("rows affected: %d\n", no)
	db.Close()

	db, err = sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, fmt.Errorf(" Error %v opening database: \n", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, fmt.Errorf(" Error %v pinging database: \n", err)
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
	return db, nil
}
