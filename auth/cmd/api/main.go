package main

import (
	"database/sql"
	"log"

	"github.com/DpodDani/auth/cmd/data"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting auth service...")
}
