package repository

import (
	"github.com/hadisjane/confessly/internal/db"
	"github.com/hadisjane/confessly/internal/errs"
	"github.com/hadisjane/confessly/internal/models"
)

// Check if confession exists
func confessionExists(confessionID int) (bool, error) {
	var exists bool
	err := db.GetDB().QueryRow("SELECT EXISTS(SELECT 1 FROM confessions WHERE id = $1)", confessionID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func CreateReport(report models.Report) error {
	// Check if confession exists
	exists, err := confessionExists(report.ConfessionID)
	if err != nil {
		return err
	}
	if !exists {
		return errs.ErrConfessionNotFound
	}

	tx, err := db.GetDB().Begin()
	if err != nil {
		return err
	}

	// Check if user has already reported this confession
	var reportExists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM reports WHERE user_id = $1 AND confession_id = $2)", report.UserID, report.ConfessionID).Scan(&reportExists)
	if err != nil {
		tx.Rollback()
		return err
	}

	if reportExists {
		tx.Rollback()
		return errs.ErrReportExists
	}

	_, err = tx.Exec("INSERT INTO reports (user_id, confession_id, reason) VALUES ($1, $2, $3)", 
		report.UserID, report.ConfessionID, report.Reason)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}