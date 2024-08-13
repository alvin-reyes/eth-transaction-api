package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"eth-transaction-api/config"
	"eth-transaction-api/models"
	"eth-transaction-api/router"
	"eth-transaction-api/seeders"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Define command-line flags
	migrateDB := flag.Bool("db:migrate", false, "Migrate the database schema")
	seedDB := flag.Bool("db:seed", false, "Seed the database with initial data")
	startServer := flag.Bool("server:start", false, "Start the HTTP server")

	flag.Parse()

	// Load configuration
	cfg := config.LoadConfig()

	// Open the SQLite database file (or create it if it doesn't exist)
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Handle the commands
	if *migrateDB {
		migrateDatabase(db)
	} else if *seedDB {
		seedDatabase(db)
	} else if *startServer {
		migrateDatabase(db)
		seedDatabase(db)
		startHTTPServer(cfg, db)
	} else {
		fmt.Println("Please specify a command: --db:migrate, --db:seed, or --server:start")
		os.Exit(1)
	}
}

// migrateDatabase handles the schema migration
func migrateDatabase(db *gorm.DB) {
	if err := db.AutoMigrate(&models.Account{}, &models.Transaction{}); err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	log.Println("database schema migrated successfully.")
}

// seedDatabase handles the database seeding
func seedDatabase(db *gorm.DB) {
	seeders.SeedAccounts(db)
}

// startHTTPServer starts the HTTP server
func startHTTPServer(cfg *config.Config, db *gorm.DB) {
	r := router.NewRouter(cfg, db)

	port := ":" + cfg.Port
	log.Printf("starting server on port %s...\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
