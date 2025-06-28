package controller

import (
	"github.com/hadisjane/confessly/internal/errs"
	"github.com/hadisjane/confessly/internal/middleware"
	"github.com/hadisjane/confessly/internal/models"
	"github.com/hadisjane/confessly/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateReport godoc
// @Summary Создание жалобы
// @Tags report
// @Accept json
// @Produce json
// @Param report body models.Report true "Report object"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /report [post]
func CreateReport(c *gin.Context) {
	// Get user ID from context
	userID := c.GetInt(middleware.UserIDCtx)
	if userID == 0 {
		HandleError(c, errs.ErrUnauthorized)
		return
	}

	// Parse request body
	var report models.Report
	if err := c.ShouldBindJSON(&report); err != nil {
		HandleError(c, err)
		return
	}

	// Validate required fields
	if report.ConfessionID <= 0 {
		HandleError(c, fmt.Errorf("confession_id is required and must be greater than 0"))
		return
	}

	if report.Reason == "" {
		HandleError(c, fmt.Errorf("reason is required"))
		return
	}

	// Set the user ID from the context
	report.UserID = userID

	// Create the report
	if err := service.CreateReport(report); err != nil {
		// Check if it's a foreign key violation
		if err.Error() == "pq: insert or update on table \"reports\" violates foreign key constraint \"reports_user_id_fkey\"" {
			HandleError(c, fmt.Errorf("user not found"))
		} else if err.Error() == "pq: insert or update on table \"reports\" violates foreign key constraint \"reports_confession_id_fkey\"" {
			HandleError(c, fmt.Errorf("confession not found"))
		} else {
			HandleError(c, err)
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Report created successfully",
	})
}
