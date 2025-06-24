package db

import "database/sql"

// Creates all necessary tables
func RunMigrations(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_img TEXT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create watchlist_items table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS watchlist_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			watchlist_id INTEGER NOT NULL,
			media_id INTEGER NOT NULL,
			watched BOOLEAN DEFAULT false,
			rating REAL,
			notes TEXT,
			added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (watchlist_id) REFERENCES watchlists(id) ON DELETE CASCADE,
			FOREIGN KEY (media_id) REFERENCES media_content(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	// Create sessions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}
	return nil
}