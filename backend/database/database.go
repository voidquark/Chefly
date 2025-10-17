package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes SQLite database connection
func InitDB(dbPath string) (*sql.DB, error) {
	// Ensure data directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	fmt.Println("✅ Database initialized successfully")
	return db, nil
}

// RunMigrations runs database migrations
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		// Create users table
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			username TEXT NOT NULL,
			is_admin INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Create recipes table
		`CREATE TABLE IF NOT EXISTS recipes (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			ingredients TEXT NOT NULL,
			steps TEXT NOT NULL,
			cooking_time INTEGER,
			difficulty TEXT,
			cuisine_type TEXT,
			meat_type TEXT,
			dietary_tags TEXT,
			is_favorite INTEGER DEFAULT 0,
			image_path TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,

		// Create indexes for recipes
		`CREATE INDEX IF NOT EXISTS idx_recipes_user_id ON recipes(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_recipes_created_at ON recipes(created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_recipes_favorite ON recipes(user_id, is_favorite)`,

		// Create recipe_filters table
		`CREATE TABLE IF NOT EXISTS recipe_filters (
			id TEXT PRIMARY KEY,
			recipe_id TEXT NOT NULL,
			meat_category TEXT,
			side_ingredients TEXT,
			country TEXT,
			dietary_preferences TEXT,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
		)`,

		// Create indexes for recipe_filters
		`CREATE INDEX IF NOT EXISTS idx_filters_recipe_id ON recipe_filters(recipe_id)`,
		`CREATE INDEX IF NOT EXISTS idx_filters_meat ON recipe_filters(meat_category)`,
		`CREATE INDEX IF NOT EXISTS idx_filters_country ON recipe_filters(country)`,

		// Create shopping_list_items table
		`CREATE TABLE IF NOT EXISTS shopping_list_items (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			recipe_id TEXT,
			recipe_title TEXT,
			ingredient_name TEXT NOT NULL,
			quantity TEXT NOT NULL,
			unit TEXT NOT NULL,
			is_checked INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE SET NULL
		)`,

		// Create indexes for shopping_list_items
		`CREATE INDEX IF NOT EXISTS idx_shopping_user_id ON shopping_list_items(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_shopping_recipe_id ON shopping_list_items(recipe_id)`,
		`CREATE INDEX IF NOT EXISTS idx_shopping_checked ON shopping_list_items(user_id, is_checked)`,

		// Migration: Add is_admin column if it doesn't exist (for existing databases)
		`ALTER TABLE users ADD COLUMN is_admin INTEGER DEFAULT 0`,

		// Migration: Add recipe_limit column to users table
		// NULL = use global limit, -1 = unlimited, 0 = blocked, >0 = custom limit
		`ALTER TABLE users ADD COLUMN recipe_limit INTEGER DEFAULT NULL`,

		// Create refresh_tokens table for JWT refresh token management
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			token TEXT NOT NULL UNIQUE,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			revoked INTEGER DEFAULT 0,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,

		// Create indexes for refresh_tokens
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at)`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			// Ignore "duplicate column" errors for ALTER TABLE migrations (migrations 13 and 14 are ALTER TABLE)
			// This allows the app to work with existing databases
			errorStr := err.Error()
			if contains(errorStr, "duplicate column") || contains(errorStr, "already exists") {
				// It's an ALTER TABLE that already ran, skip it
				continue
			}
			return fmt.Errorf("migration %d failed: %w", i, err)
		}
	}

	fmt.Println("✅ Database migrations completed successfully")
	return nil
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
