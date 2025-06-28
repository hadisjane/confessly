package repository

import (
	"Confessly/internal/db"
	"Confessly/internal/errs"
	"Confessly/internal/models"
)

func GetUserByUsernameAndPassword(username string, password string) (models.User, error) {
	var user models.User
	err := db.GetDB().Get(&user, `
		SELECT id, username, role, created_at
		FROM users 
		WHERE username = $1 AND password = $2`,
		username, password)

	if err != nil {
		return models.User{}, translateError(err)
	}

	return user, nil
}

func GetUserByUsername(username string) (user models.User, err error) {
	err = db.GetDB().Get(&user, `SELECT id, 
					   username,
					   role,
					   password,
					   created_at
				FROM users WHERE username = $1`, username)
	if err != nil {
		return models.User{}, translateError(err)
	}

	return user, nil
}

func CreateUser(user models.UserRegister) error {
	_, err := db.GetDB().Exec(`
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)`,
		user.Username,
		user.Email,
		user.Password)
	return err
}

// UpdateUserBannedStatus updates the banned status of a user
func UpdateUserBannedStatus(id int, banned bool) error {
	result, err := db.GetDB().Exec(`
		UPDATE users 
		SET banned = $1 
		WHERE id = $2`, banned, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func IsUserBanned(id int) (bool, error) {
	var banned bool
	err := db.GetDB().Get(&banned, "SELECT banned FROM users WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	return banned, nil
}