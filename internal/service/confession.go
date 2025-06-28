package service

import (
	"github.com/hadisjane/confessly/internal/models"
	"github.com/hadisjane/confessly/internal/repository"
	"errors"
)

// CreateConfession creates a new confession
func CreateConfession(confession models.Confession) error {
	// Validate that either UserID or GuestUUID is set, but not both
	if (confession.UserID == nil && confession.GuestUUID == nil) || 
	   (confession.UserID != nil && confession.GuestUUID != nil) {
		return errors.New("confession must have either user ID or guest UUID")
	}
	
	return repository.CreateConfession(confession)
}
// GetConfessions retrieves all confessions
func GetAllConfessions() ([]models.Confession, error) {
	return repository.GetAllConfessions()
}

// GetConfession retrieves a single confession by ID
func GetConfession(id int) (models.Confession, error) {
	return repository.GetConfession(id)
}

// UpdateConfession updates an existing confession
func UpdateConfession(id int, confession models.Confession) error {
	return repository.UpdateConfession(id, confession)
}

// DeleteConfession deletes a confession by ID
func DeleteConfession(id int) error {
	return repository.DeleteConfession(id)
}

// SearchConfessionsByTitle searches confessions by title
func SearchConfessionsByTitle(title string) ([]models.Confession, error) {
	return repository.SearchConfessionsByTitle(title)
}
