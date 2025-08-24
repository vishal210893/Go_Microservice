// Package main provides a database seeding utility for populating the database
// with test data for development and testing environments.
//
// The seeder connects to a PostgreSQL database and creates sample users, posts,
// and comments using predefined test data.
//
// Usage:
//
//	go run cmd/migrate/seed/main.go
//
// Environment Variables:
//
//	DB_ADDR - PostgreSQL connection string (optional, has default)
package main

import (
	"Go-Microservice/internal/db"
	"Go-Microservice/internal/env"
	"Go-Microservice/internal/repo"
	"log"
)

// Default database configuration constants
const (
	// defaultDBAddr is the default PostgreSQL connection string
	defaultDBAddr = "postgres://avnadmin:AVNS_LT5DsEKUPKfrHSHZHyB@pg-1d9d15dc-vishal210893-5985.h.aivencloud.com:28832/defaultdb?sslmode=require"

	// Database connection pool configuration
	maxOpenConns = 3
	maxIdleConns = 3
	maxIdleTime  = "15m"
)

// main initializes the database connection and runs the seeding process.
// It retrieves the database address from environment variables, establishes
// a connection with configured pool settings, and populates the database
// with test data.
//
// The process will terminate with a fatal error if:
//   - Database connection fails
//   - Repository initialization fails
//
// On successful completion, the database will contain sample users, posts, and comments.
func main() {
	// Get database connection string from environment or use default
	dbAddr := env.GetString("DB_ADDR", defaultDBAddr)

	// Establish database connection with pool configuration
	conn, err := db.New(dbAddr, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Initialize repository with database connection
	repository, err := repo.NewPostgresRepo(conn)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Run database seeding
	log.Println("Starting database seeding process...")
	db.Seed(*repository)
	log.Println("Database seeding completed successfully")
}
