package service

import (
	"Confessly/internal/models"
	"Confessly/internal/repository"
)

func CreateGuestUser(guestUser models.GuestUser) error {
	return repository.CreateGuestUser(guestUser)
}

func GetGuestUsers() ([]models.GuestUser, error) {
	return repository.GetGuestUsers()
}

func GetGuestUser(uuid string) (models.GuestUser, error) {
	return repository.GetGuestUser(uuid)
}

func IsGuestBanned(uuid string) (bool, error) {
	return repository.IsGuestBanned(uuid)
}