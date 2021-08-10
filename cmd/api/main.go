package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn      string
		dbName   string
		username string
		password string
		host     string
		port     string
	}
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production)")
	flag.StringVar(&cfg.db.host, "db-host", os.Getenv("POSTGRES_HOST"), "Postgres host adress")
	flag.StringVar(&cfg.db.port, "db-port", os.Getenv("POSTGRES_PORT"), "Postgres port adress")
	flag.StringVar(&cfg.db.dbName, "db-name", os.Getenv("POSTGRES_DB"), "Postgres database name")
	flag.StringVar(&cfg.db.username, "db-username", os.Getenv("POSTGRES_USER"), "Postgres username")
	flag.StringVar(&cfg.db.password, "db-password", os.Getenv("POSTGRES_PASSWORD"), "Postgres password secret")
	flag.StringVar(&cfg.db.dsn, "dsn", "postgresql://"+cfg.db.username+":"+cfg.db.password+"@"+cfg.db.host+":"+cfg.db.port+"/"+cfg.db.dbName+"?sslmode=disable", "Postgres connection string")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port", cfg.port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
