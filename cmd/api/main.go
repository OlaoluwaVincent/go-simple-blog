package main

import (
	"log"

	"github.com/olaoluwavincent/full-course/internal/db"
	"github.com/olaoluwavincent/full-course/internal/env"
	"github.com/olaoluwavincent/full-course/internal/store"
)

func main() {
	// Initialize configuration from environment variables
	cfg := config{
		addr: env.GetEnvString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetEnvString("DB_URL", "postgres://your_db_username:your_password@host:your_port/your_db_name?sslmode=disable"),
			maxOpenConns: env.GetEnvInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetEnvInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetEnvString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	// Initialize database connection
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Database connection established")

	// Initialize storage layer
	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	log.Fatal(app.run(app.mount()))
}
