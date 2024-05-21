package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DpodDani/auth/cmd/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

type Config struct {
	Repo   data.Repository
	Client *http.Client
}

func main() {
	dsn := os.Getenv("DSN")

	log.Println("Starting auth service...")

	// connect to DB
	db := connectToDB(dsn)

	if db == nil {
		log.Panic("Can't connect to database! 🔥")
	}

	app := Config{
		Client: &http.Client{},
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB(dsn string) *sql.DB {
	counts := 0
	backOffTime := 2 * time.Second

	for {
		db, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return db
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Printf("Backing off for %d seconds...\n", backOffTime)
		time.Sleep(backOffTime)
	}
}

func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewPostgresRepository(conn)
	app.Repo = db
}
