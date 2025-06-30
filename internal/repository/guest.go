package repository

import (
	"time"

	"github.com/hadisjane/confessly/internal/db"
	"github.com/hadisjane/confessly/internal/models"
)

func CreateGuestUser(guestUser models.GuestUser) error {
	tx, err := db.GetDB().Begin()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO guest_users (uuid, banned, created_at)
		VALUES ($1, $2, $3)
		RETURNING uuid
	`

	err = tx.QueryRow(
		query,
		guestUser.UUID,
		guestUser.Banned,
		time.Now(),
	).Scan(&guestUser.UUID)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetGuestUser(uuid string) (models.GuestUser, error) {
	var guestUser models.GuestUser

	err := db.GetDB().Get(&guestUser, "SELECT uuid, banned, created_at FROM guest_users WHERE uuid = $1", uuid)
	if err != nil {
		return models.GuestUser{}, err
	}
	return guestUser, nil
}


func IsGuestBanned(uuid string) (bool, error) {
	var guestUser models.GuestUser

	err := db.GetDB().Get(&guestUser, "SELECT banned FROM guest_users WHERE uuid = $1", uuid)
	if err != nil {
		return false, err
	}
	return guestUser.Banned, nil
}

func GetGuestUsers() ([]models.GuestUser, error) {
	var guestUsers []models.GuestUser
	err := db.GetDB().Select(&guestUsers, "SELECT uuid, banned, created_at FROM guest_users")
	if err != nil {
		return nil, err
	}
	return guestUsers, nil
}