package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"deepphotos/backend/internal/config"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	dir := filepath.Dir(cfg.DBPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;`); err != nil {
		log.Printf("Warning: failed to set WAL mode: %v", err)
	}

	if err := migrateSchema(db); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	if err := seedDefaultAdmin(db, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		log.Printf("Warning: admin seeding check: %v", err)
	}

	log.Printf("SQLite database initialized successfully at %s", cfg.DBPath)
	return db, nil
}

func migrateSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'Viewer',
		avatar TEXT,
		status TEXT NOT NULL DEFAULT 'Active',
		last_login TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS photos (
		id TEXT PRIMARY KEY,
		filename TEXT NOT NULL,
		object_key TEXT NOT NULL,
		thumbnail_key TEXT NOT NULL,
		mime_type TEXT NOT NULL,
		file_type TEXT NOT NULL DEFAULT 'image',
		size INTEGER NOT NULL,
		width INTEGER DEFAULT 0,
		height INTEGER DEFAULT 0,
		exif_model TEXT,
		taken_at TEXT,
		uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		latitude REAL DEFAULT 0,
		longitude REAL DEFAULT 0,
		hash TEXT,
		is_favorite INTEGER DEFAULT 0,
		is_deleted INTEGER DEFAULT 0,
		locked_folder_id TEXT,
		title TEXT
	);

	CREATE TABLE IF NOT EXISTS locked_folders (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		passcode_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS albums (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		cover_url TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS album_photos (
		album_id TEXT NOT NULL,
		photo_id TEXT NOT NULL,
		PRIMARY KEY (album_id, photo_id),
		FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE,
		FOREIGN KEY (photo_id) REFERENCES photos(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS login_logs (
		id TEXT PRIMARY KEY,
		user_email TEXT NOT NULL,
		timestamp TEXT NOT NULL,
		ip TEXT NOT NULL,
		device TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	// Safe alter columns for existing databases
	db.Exec(`ALTER TABLE photos ADD COLUMN file_type TEXT NOT NULL DEFAULT 'image';`)
	db.Exec(`ALTER TABLE photos ADD COLUMN exif_model TEXT;`)
	db.Exec(`ALTER TABLE photos ADD COLUMN locked_folder_id TEXT;`)

	return nil
}

func seedDefaultAdmin(db *sql.DB, email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		query := `
		INSERT INTO users (id, name, email, password, role, avatar, status, last_login)
		VALUES ('usr_admin_1', 'Deepak (Admin)', ?, ?, 'Administrator', 
		'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150&auto=format&fit=crop&q=80', 'Active', 'Just now')
		`
		_, err = db.Exec(query, email, string(hashed))
		if err != nil {
			return err
		}
		log.Printf("Seeded default admin account: %s", email)
	} else {
		query := `UPDATE users SET password = ? WHERE email = ?`
		_, err = db.Exec(query, string(hashed), email)
		if err != nil {
			return err
		}
		log.Printf("Updated admin account credentials for: %s", email)
	}

	return nil
}
