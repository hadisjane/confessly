package service

import (
	"github.com/hadisjane/confessly/internal/models"
	"github.com/hadisjane/confessly/internal/repository"
)

func CreateReport(report models.Report) error {
	return repository.CreateReport(report)
}
	