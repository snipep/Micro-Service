package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/snipep/authentication-service/data"
)

const webPort = "80"

var counts int64

// Config holds the application's configuration, including database and models
type Config struct {
	DB     	*sql.DB
	Models 	data.Models
}

func main() {
	log.Println("Starting authentication service")

	// Attempt to connect to the database
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// Set up application configuration with database connection and models
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	// Create and configure the HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// Start the server and listen for incoming requests
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// openDB opens a new database connection using the provided DSN
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Verify the database connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// connectToDB attempts to establish a connection to the database, retrying if necessary
func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		// Give up after 10 unsuccessful attempts
		if counts > 10 {
			log.Println(err)
			return nil
		}

		// Back off before retrying
		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
