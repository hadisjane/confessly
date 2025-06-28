package repository

import (
	"Confessly/internal/db"
	"Confessly/internal/errs"
	"Confessly/internal/models"
	"database/sql"
	"fmt"
)
	
func GetReports() []models.Report {
	reports := make([]models.Report, 0) // Initialize empty slice
	err := db.GetDB().Select(&reports, "SELECT * FROM reports")
	if err != nil {
		return reports // Return empty slice instead of nil
	}	
	return reports
}

func GetUsers() []models.User {
	users := make([]models.User, 0) // Initialize empty slice
	err := db.GetDB().Select(&users, "SELECT id, username, email, role, banned, created_at FROM users")
	if err != nil {
		return users // Return empty slice instead of nil
	}	
	return users
}

// GetUserByID retrieves a user by their ID
func GetUserByID(id int) (models.User, error) {
	var user models.User
	err := db.GetDB().Get(&user, `
		SELECT id, username, email, role, banned, created_at
		FROM users 
		WHERE id = $1`, id)

	if err == sql.ErrNoRows {
		return models.User{}, errs.ErrNotFound
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func BanUser(userID int) error {
	_, err := db.GetDB().Exec("UPDATE users SET banned = true WHERE id = $1", userID)
	if err != nil {
		return err
	}
	return nil
}

func UnbanUser(userID int) error {
	_, err := db.GetDB().Exec("UPDATE users SET banned = false WHERE id = $1", userID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteConfessionByAdmin(confessionID int) error {
	tx, err := db.GetDB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// First delete all reports associated with this confession
	_, err = tx.Exec("DELETE FROM reports WHERE confession_id = $1", confessionID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete reports: %w", err)
	}

	// Then delete the confession
	_, err = tx.Exec("DELETE FROM confessions WHERE id = $1", confessionID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete confession: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func BanGuestUser(uuid string) error {
	_, err := db.GetDB().Exec("UPDATE guest_users SET banned = true WHERE uuid = $1", uuid)
	if err != nil {
		return err
	}
	return nil
}

func UnbanGuestUser(uuid string) error {
	_, err := db.GetDB().Exec("UPDATE guest_users SET banned = false WHERE uuid = $1", uuid)
	if err != nil {
		return err
	}
	return nil
}

func UpdateReport(reportID int, updateReq models.UpdateReport) error {
	_, err := db.GetDB().Exec("UPDATE reports SET status = $1 WHERE id = $2", updateReq.Status, reportID)
	if err != nil {
		return err
	}
	return nil
}

func GetReport(reportID int) (models.Report, error) {
	var report models.Report
	err := db.GetDB().Get(&report, "SELECT * FROM reports WHERE id = $1", reportID)
	if err == sql.ErrNoRows {
		return models.Report{}, errs.ErrNotFound
	}
	if err != nil {
		return models.Report{}, err
	}
	return report, nil
}
