package service

import (
	"Confessly/internal/models"
	"Confessly/internal/repository"
)

func CreateReport(report models.Report) error {
	return repository.CreateReport(report)
}
	