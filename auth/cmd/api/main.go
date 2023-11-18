package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/DpodDani/auth/cmd/data"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting auth service...")

	// TODO: Setup DB connection

	app := Config{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
