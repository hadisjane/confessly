package service

import (
	"Confessly/internal/errs"
	"Confessly/internal/models"
	"Confessly/internal/repository"
	"Confessly/utils"
	"errors"
)

func CreateUser(u models.UserRegister) error {
	_, err := repository.GetUserByUsername(u.Username)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			// User doesn't exist, we can proceed with creation
		} else {
			return err
		}
	} else {
		return errs.ErrUserAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	if err = repository.CreateUser(u); err != nil {
		return err
	}

	return nil
}

func GetUserByUsernameAndPassword(username string, password string) (models.User, error) {
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return models.User{}, errs.ErrIncorrectUsernameOrPassword
		}
		return models.User{}, err
	}

	// Check if user is banned
	if user.Banned {
		return models.User{}, errs.ErrUserBanned
	}

	err = utils.VerifyPassword(user.Password, password)
	if err != nil {
		return models.User{}, errs.ErrIncorrectUsernameOrPassword
	}

	return user, nil
}

// GetUser retrieves a user by ID
func GetUser(id int) (models.User, error) {
	user, err := repository.GetUserByID(id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func IsUserBanned(id int) (bool, error) {
	return repository.IsUserBanned(id)
}
