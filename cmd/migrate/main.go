package main

import (
	"database/sql"
	"flag"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/w0ikid/highload-auth-go/pkg/config"
)

func main() {
	var dir string
	var cmd string

	flag.StringVar(&dir, "dir", "migrations", "directory with migration files")
	flag.StringVar(&cmd, "cmd", "up", "goose command (up, down, status)")
	flag.Parse()

	cfg := config.Load()
	dsn := cfg.Postgres.DSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v\n", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v\n", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v\n", err)
	}

	log.Printf("running goose %s...", cmd)

	switch cmd {
	case "up":
		if err := goose.Up(db, dir); err != nil {
			log.Fatalf("goose up failed: %v\n", err)
		}
	case "down":
		if err := goose.Down(db, dir); err != nil {
			log.Fatalf("goose down failed: %v\n", err)
		}
	case "status":
		if err := goose.Status(db, dir); err != nil {
			log.Fatalf("goose status failed: %v\n", err)
		}
	default:
		log.Fatalf("unknown command: %s. expected up, down or status", cmd)
	}

	log.Printf("goose %s completed successfully", cmd)
}
