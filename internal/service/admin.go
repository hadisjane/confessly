package service

import (
	"Confessly/internal/models"
	"Confessly/internal/repository"
)

func GetReports() []models.Report {
	return repository.GetReports()
}	

func GetUsers() []models.User {
	return repository.GetUsers()
}

func GetUserByID(id int) (models.User, error) {
	return repository.GetUserByID(id)
}

// BanUser bans a user by ID
func BanUser(id int) error {
	// First check if user exists
	_, err := repository.GetUserByID(id)
	if err != nil {
		return err
	}

	// Update user's banned status
	return repository.UpdateUserBannedStatus(id, true)
}

// UnbanUser unbans a user by ID
func UnbanUser(id int) error {
	// First check if user exists
	_, err := repository.GetUserByID(id)
	if err != nil {
		return err
	}

	// Update user's banned status
	return repository.UpdateUserBannedStatus(id, false)
}

func DeleteConfessionByAdmin(confessionID int) error {
	return repository.DeleteConfessionByAdmin(confessionID)
}

func BanGuestUser(uuid string) error {
	return repository.BanGuestUser(uuid)
}

func UnbanGuestUser(uuid string) error {
	return repository.UnbanGuestUser(uuid)
}

func UpdateReport(reportID int, updateReq models.UpdateReport) error {
	return repository.UpdateReport(reportID, updateReq)
}

func GetReport(reportID int) (models.Report, error) {
	return repository.GetReport(reportID)
}
	