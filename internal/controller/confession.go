package controller

import (
	"net/http"
	"strconv"

	"github.com/hadisjane/confessly/internal/errs"
	"github.com/hadisjane/confessly/internal/middleware"
	"github.com/hadisjane/confessly/internal/models"
	"github.com/hadisjane/confessly/internal/service"

	"github.com/gin-gonic/gin"
)

// Ping godoc
// @Summary Проверка работоспособности сервера
// @Tags health
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router / [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Confessly server up and running",
	})
}

type CreateConfessionRequest struct {
	Title string `json:"title" binding:"required"`
	Text  string `json:"text" binding:"required"`
	Anon  bool   `json:"anon"`
}

// CreateConfession godoc
// @Summary Создание конфесии
// @Tags confession
// @Accept json
// @Produce json
// @Param confession body CreateConfessionRequest true "Confession object"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /confessions [post]
func CreateConfession(c *gin.Context) {
	var req CreateConfessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	userID := c.GetInt(middleware.UserIDCtx)
	username, userExists := c.Get(middleware.UsernameCtx)

	var confession models.Confession

	if !userExists || userID == 0 {
		// Handle guest user
		guestUUID, exists := c.Get(middleware.GuestUUIDCtx)
		if !exists {
			HandleError(c, errs.ErrUnauthorized)
			return
		}

		guestUUIDStr, ok := guestUUID.(string)
		if !ok {
			HandleError(c, errs.ErrInternalServer)
			return
		}

		// Get guest user from database
		_, err := service.GetGuestUser(guestUUIDStr)
		if err != nil {
			HandleError(c, err)
			return
		}

		confession = models.Confession{
			GuestUUID: &guestUUIDStr,
			Username:  "Guest_" + guestUUIDStr[:8],
			Title:     req.Title,
			Text:      req.Text,
			Anon:      true,
		}
	} else {
		// Handle authenticated user
		usernameStr, ok := username.(string)
		if !ok {
			HandleError(c, errs.ErrUnauthorized)
			return
		}

		confession = models.Confession{
			UserID:   &userID,
			Username: usernameStr,
			Title:    req.Title,
			Text:     req.Text,
			Anon:     req.Anon,
		}
	}

	if err := service.CreateConfession(confession); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Confession created successfully",
	})
}

// GetAllConfessions godoc
// @Summary Получение всех конфесий
// @Tags confession
// @Produce json
// @Success 200 {object} []models.Confession
// @Failure 500 {object} map[string]string
// @Router /confessions [get]
func GetAllConfessions(c *gin.Context) {
	userRole, _ := c.Get(middleware.RoleCtx)
	userRoleStr, _ := userRole.(string)

	confessions, err := service.GetAllConfessions()
	if err != nil {
		HandleError(c, err)
		return
	}

	// Hide sensitive info for non-admin users
	if userRoleStr != "admin" {
		for i := range confessions {
			if confessions[i].Anon {
				confessions[i].UserID = nil
				confessions[i].GuestUUID = nil
				confessions[i].Username = ""
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"confessions": confessions,
	})
}

// GetConfession godoc
// @Summary Получение конфесии по ID
// @Tags confession
// @Produce json
// @Param id path int true "Confession ID"
// @Success 200 {object} models.Confession
// @Failure 404 {object} map[string]string
// @Router /confessions/{id} [get]
func GetConfession(c *gin.Context) {
	userRole, _ := c.Get(middleware.RoleCtx)
	userRoleStr, _ := userRole.(string)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}

	confession, err := service.GetConfession(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Hide sensitive info for non-admin users
	if confession.Anon && userRoleStr != "admin" {
		confession.UserID = nil
		confession.GuestUUID = nil
		confession.Username = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"confession": confession,
	})
}

// UpdateConfessionRequest defines the structure for updating a confession
type UpdateConfessionRequest struct {
	Title *string `json:"title,omitempty"`
	Text  *string `json:"text,omitempty"`
	Anon  *bool   `json:"anon,omitempty"`
}

// UpdateConfession godoc
// @Summary Обновление конфесии
// @Tags confession
// @Accept json
// @Produce json
// @Param id path int true "Confession ID"
// @Param confession body UpdateConfessionRequest true "Confession object"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /confessions/{id} [put]
func UpdateConfession(c *gin.Context) {
	userID := c.GetInt(middleware.UserIDCtx)
	if userID == 0 {
		HandleError(c, errs.ErrUnauthorized)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}

	// Get the confession first to check ownership
	existingConfession, err := service.GetConfession(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Check if the requesting user is the owner
	if existingConfession.UserID == nil || *existingConfession.UserID != userID {
		HandleError(c, errs.ErrForbidden)
		return
	}

	var updateReq UpdateConfessionRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		HandleError(c, err)
		return
	}

	// Create an update confession with existing values
	updatedConfession := models.Confession{
		ID:       id,
		UserID:   &userID,
		Username: existingConfession.Username,
		Title:    existingConfession.Title,
		Text:     existingConfession.Text,
		Anon:     existingConfession.Anon,
	}

	// Update only the fields that were provided in the request
	if updateReq.Title != nil {
		updatedConfession.Title = *updateReq.Title
	}
	if updateReq.Text != nil {
		updatedConfession.Text = *updateReq.Text
	}
	if updateReq.Anon != nil {
		updatedConfession.Anon = *updateReq.Anon
	}

	if err := service.UpdateConfession(id, updatedConfession); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Confession updated successfully",
	})
}

// DeleteConfession godoc
// @Summary Удаление конфесии
// @Tags confession
// @Param id path int true "Confession ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /confessions/{id} [delete]
func DeleteConfession(c *gin.Context) {
	userID := c.GetInt(middleware.UserIDCtx)
	if userID == 0 {
		HandleError(c, errs.ErrUnauthorized)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}

	// Get the confession first to check ownership
	confession, err := service.GetConfession(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Check if the requesting user is the owner
	if confession.UserID == nil || *confession.UserID != userID {
		HandleError(c, errs.ErrForbiddenDelete)
		return
	}

	if err := service.DeleteConfession(id); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Confession deleted successfully",
	})
}

// SearchConfessions godoc
// @Summary Поиск конфесий по названию
// @Tags confession
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} []models.Confession
// @Failure 500 {object} map[string]string
// @Router /confessions/search [get]
func SearchConfessions(c *gin.Context) {
	query := c.Query("q")

	if query == "" {
		confessions, err := service.GetAllConfessions()
		if err != nil {
			HandleError(c, err)
			return
		}

		// Hide sensitive info for non-admin users
		userRole, _ := c.Get(middleware.RoleCtx)
		userRoleStr, _ := userRole.(string)

		if userRoleStr != "admin" {
			for i := range confessions {
				if confessions[i].Anon {
					confessions[i].UserID = nil
					confessions[i].GuestUUID = nil
					confessions[i].Username = ""
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"confessions": confessions,
		})
		return
	}

	confessions, err := service.SearchConfessionsByTitle(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	if confessions == nil {
		confessions = []models.Confession{}
	}

	// Hide sensitive info for non-admin users
	userRole, _ := c.Get(middleware.RoleCtx)
	userRoleStr, _ := userRole.(string)

	if userRoleStr != "admin" {
		for i := range confessions {
			if confessions[i].Anon {
				confessions[i].UserID = nil
				confessions[i].GuestUUID = nil
				confessions[i].Username = ""
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"confessions": confessions,
	})
}
