package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB represents the database connection
type DB struct {
	*sqlx.DB
}

// Initialize creates a new database connection
func Initialize(dataSourceName string) (*DB, error) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	
	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * 60 * time.Second)
	
	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}
	
	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, err
	}
	
	return &DB{db}, nil
}

// createTables ensures all required tables exist
func createTables(db *sqlx.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(20) NOT NULL DEFAULT 'user',
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}
	
	// Create categories table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) UNIQUE NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}
	
	// Create tasks table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			description TEXT,
			user_id INT REFERENCES users(id) ON DELETE CASCADE,
			category_id INT REFERENCES categories(id) ON DELETE SET NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			due_date TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}
	
	// Create audit logs table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS audit_logs (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id) ON DELETE SET NULL,
			action VARCHAR(50) NOT NULL,
			entity_type VARCHAR(50) NOT NULL,
			entity_id INT,
			details JSONB,
			ip_address VARCHAR(50),
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	
	return err
}