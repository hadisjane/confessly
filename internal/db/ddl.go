package db

import (
	"fmt"
	"log"
)



func InitMigrations() error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	// Сначала создаем таблицу пользователей
	usersTable := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		email VARCHAR(255) NOT NULL UNIQUE,
		role VARCHAR(255) NOT NULL DEFAULT 'user',
		password VARCHAR(255) NOT NULL,
		banned BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`
	log.Println("Creating users table if not exists...")

	// Создаем таблицу пользователей
	if _, err := db.Exec(usersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	guestUsersTable := `CREATE TABLE IF NOT EXISTS guest_users (
		uuid UUID PRIMARY KEY,
		banned BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`
	log.Println("Creating guest users table if not exists...")

	// Создаем таблицу гостей
	if _, err := db.Exec(guestUsersTable); err != nil {
		return fmt.Errorf("failed to create guest users table: %w", err)
	}

	// Create confessions table with proper foreign key to users
	confessionsTable := `CREATE TABLE IF NOT EXISTS confessions (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
		guest_uuid UUID REFERENCES guest_users(uuid) ON DELETE SET NULL,
		username VARCHAR(255) NOT NULL,
		title VARCHAR(100) NOT NULL DEFAULT 'Untitled',
		text TEXT NOT NULL,
		anon BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		
		-- Ensure either user_id or guest_uuid is set
		CONSTRAINT chk_user_or_guest CHECK (
			(user_id IS NOT NULL AND guest_uuid IS NULL) OR
			(user_id IS NULL AND guest_uuid IS NOT NULL)
		)
	)`

	log.Println("Creating confession table if not exists...")
	// Создаем таблицу признаний
	if _, err := db.Exec(confessionsTable); err != nil {
		return fmt.Errorf("failed to create confession table: %w", err)
	}

	// Create reports table with proper foreign key to users
	reportsTable := `CREATE TABLE IF NOT EXISTS reports (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) NOT NULL,
		confession_id INTEGER REFERENCES confessions(id) NOT NULL,
		reason TEXT NOT NULL,
		status VARCHAR(255) NOT NULL DEFAULT 'pending',
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	log.Println("Creating reports table if not exists...")
	// Создаем таблицу жалоб
	if _, err := db.Exec(reportsTable); err != nil {
		return fmt.Errorf("failed to create reports table: %w", err)
	}

	log.Println("Database migrations completed successfully")

	return nil
}