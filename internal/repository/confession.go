package repository

import (
	"Confessly/internal/db"
	"Confessly/internal/errs"
	"Confessly/internal/models"
	"database/sql"
	"fmt"
	"time"
)

// CreateConfession creates a new confession in the database
func CreateConfession(confession models.Confession) error {
	tx, err := db.GetDB().Begin()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO confessions (
			user_id, 
			guest_uuid, 
			username, 
			title, 
			text, 
			anon, 
			created_at, 
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var userID interface{}
	var guestUUID interface{}

	if confession.UserID != nil {
		userID = *confession.UserID
	} else {
		userID = nil
	}

	if confession.GuestUUID != nil {
		guestUUID = *confession.GuestUUID
	} else {
		guestUUID = nil
	}

	now := time.Now()
	err = tx.QueryRow(
		query,
		userID,
		guestUUID,
		confession.Username,
		confession.Title,
		confession.Text,
		confession.Anon,
		now,
		now,
	).Scan(&confession.ID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create confession: %w", err)
	}

	return tx.Commit()
}

// GetAllConfessions retrieves all confessions from the database
func GetAllConfessions() ([]models.Confession, error) {
	var confessions []models.Confession

	query := `
		SELECT 
			id, 
			user_id, 
			guest_uuid, 
			username, 
			title, 
			text, 
			anon, 
			created_at, 
			updated_at
		FROM confessions
		ORDER BY created_at DESC
	`

	err := db.GetDB().Select(&confessions, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.Confession{}, nil
		}
		return nil, err
	}

	return confessions, nil
}

// GetConfession retrieves a single confession by ID
func GetConfession(id int) (models.Confession, error) {
	var confession models.Confession

	query := `
		SELECT 
			id, 
			user_id, 
			guest_uuid, 
			username, 
			title, 
			text, 
			anon, 
			created_at, 
			updated_at
		FROM confessions
		WHERE id = $1
	`

	err := db.GetDB().Get(&confession, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Confession{}, errs.ErrNotFound
		}
		return models.Confession{}, err
	}

	return confession, nil
}

// UpdateConfession updates an existing confession
func UpdateConfession(id int, confession models.Confession) error {
	tx, err := db.GetDB().Begin()
	if err != nil {
		return err
	}

	query := `
		UPDATE confessions
		SET 
			title = $1, 
			text = $2, 
			anon = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING id
	`

	var updatedID int
	err = tx.QueryRow(
		query,
		confession.Title,
		confession.Text,
		confession.Anon,
		time.Now(),
		id,
	).Scan(&updatedID)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return errs.ErrNotFound
		}
		return err
	}

	return tx.Commit()
}

// DeleteConfession deletes a confession by ID
func DeleteConfession(id int) error {
	tx, err := db.GetDB().Begin()
	if err != nil {
		return err
	}

	query := `
		DELETE FROM confessions
		WHERE id = $1
		RETURNING id
	`

	var deletedID int
	err = tx.QueryRow(query, id).Scan(&deletedID)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return errs.ErrNotFound
		}
		return err
	}

	return tx.Commit()
}

// SearchConfessionsByTitle searches confessions by title
func SearchConfessionsByTitle(searchQuery string) ([]models.Confession, error) {
	query := `
		SELECT 
			id, 
			user_id, 
			guest_uuid, 
			username, 
			title, 
			text, 
			anon, 
			created_at, 
			updated_at
		FROM confessions
		WHERE title ILIKE $1
		ORDER BY created_at DESC
		LIMIT 100
	`

	var confessions []models.Confession
	searchTerm := "%" + searchQuery + "%"
	

	err := db.GetDB().Select(&confessions, query, searchTerm)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.Confession{}, nil
		}
		return nil, fmt.Errorf("failed to search confessions: %w", err)
	}

	return confessions, nil
}