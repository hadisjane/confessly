package controller

import (
	"github.com/hadisjane/confessly/internal/errs"
	"github.com/hadisjane/confessly/internal/middleware"
	"github.com/hadisjane/confessly/internal/models"
	"github.com/hadisjane/confessly/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetReports godoc
// @Summary Получение всех жалоб (только для администраторов)
// @Tags admin
// @Produce json
// @Success 200 {object} []models.Report
// @Failure 500 {object} map[string]string
// @Router /admin/reports [get]
func GetReports(c *gin.Context) {
	reports := service.GetReports()
	if reports == nil {
		HandleError(c, errs.ErrNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
	})		
}

// GetUsers godoc
// @Summary Получение всех пользователей (только для администраторов)
// @Tags admin
// @Produce json
// @Success 200 {object} []models.User
// @Failure 500 {object} map[string]string
// @Router /admin/users [get]
func GetUsers(c *gin.Context) {
	users := service.GetUsers()
	if users == nil {
		HandleError(c, errs.ErrNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})		
}

// GetUserByID godoc
// @Summary Получение пользователя по ID (только для администраторов)
// @Tags admin
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users/{id} [get]
func GetUserByID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID <= 0 {
		HandleError(c, errs.ErrInvalidId)
		return
	}

	// First check if user exists
	user, err := service.GetUserByID(userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// DeleteConfessionByAdmin godoc
// @Summary Удаление конфесии по ID (только для администраторов)
// @Tags admin
// @Param id path int true "Confession ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/confessions/{id} [delete]
func DeleteConfessionByAdmin(c *gin.Context) {
	// Get user role from context
	userRole, exists := c.Get(middleware.RoleCtx)
	if !exists || userRole != "admin" {
		HandleError(c, errs.ErrForbidden)
		return
	}

	confessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil || confessionID <= 0 {
		HandleError(c, errs.ErrInvalidId)
		return
	}

	// First check if confession exists
	_, err = service.GetConfession(confessionID)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Delete the confession
	if err := service.DeleteConfessionByAdmin(confessionID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Confession deleted successfully by admin",
	})
}

// BanUser godoc
// @Summary Бан пользователя (только для администраторов)
// @Tags admin
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users/{id}/ban [post]
func BanUser(c *gin.Context) {
	// Get user role from context
	userRole, exists := c.Get(middleware.RoleCtx)
	if !exists || userRole != "admin" {
		HandleError(c, errs.ErrForbidden)
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID <= 0 {
		HandleError(c, errs.ErrInvalidId)
		return
	}

	// First check if user exists and get current ban status
	user, err := service.GetUser(userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Check if user is already banned
	if user.Banned {
		HandleError(c, errors.New("user is already banned"))
		return
	}

	// Ban the user
	if err := service.BanUser(userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User banned successfully",
	})
}

// UnbanUser godoc
// @Summary Разбан пользователя (только для администраторов)
// @Tags admin
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users/{id}/unban [post]
func UnbanUser(c *gin.Context) {
	// Get user role from context
	userRole, exists := c.Get(middleware.RoleCtx)
	if !exists || userRole != "admin" {
		HandleError(c, errs.ErrForbidden)
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID <= 0 {
		HandleError(c, errs.ErrInvalidId)
		return
	}

	// First check if user exists and get current ban status
	user, err := service.GetUser(userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Check if user is already unbanned
	if !user.Banned {
		HandleError(c, errors.New("user is not banned"))
		return
	}

	// Unban the user
	if err := service.UnbanUser(userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User unbanned successfully",
	})
}

// BanGuestUser godoc
// @Summary Бан гостевого пользователя (только для администраторов)
// @Tags admin
// @Param uuid path string true "Guest UUID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/guest/{uuid}/ban [post]
func BanGuestUser(c *gin.Context) {
	// Get user role from context
	userRole, exists := c.Get(middleware.RoleCtx)
	if !exists || userRole != "admin" {
		HandleError(c, errs.ErrForbidden)
		return
	}

	uuid := c.Param("uuid")

	// First check if guest user exists
	_, err := service.GetGuestUser(uuid)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Check if guest user is already banned
	isBanned, err := service.IsGuestBanned(uuid)
	if err != nil {
		HandleError(c, err)
		return
	}

	if isBanned {
		HandleError(c, errors.New("guest user is already banned"))
		return
	}

	// Ban the guest user
	if err := service.BanGuestUser(uuid); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Guest user banned successfully",
	})
}

// UnbanGuestUser godoc
// @Summary Разбан гостевого пользователя (только для администраторов)
// @Tags admin
// @Param uuid path string true "Guest UUID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/guest/{uuid}/unban [post]
func UnbanGuestUser(c *gin.Context) {
	// Get user role from context
	userRole, exists := c.Get(middleware.RoleCtx)
	if !exists || userRole != "admin" {
		HandleError(c, errs.ErrForbidden)
		return
	}

	uuid := c.Param("uuid")

	// First check if guest user exists
	_, err := service.GetGuestUser(uuid)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Check if guest user is already unbanned
	isBanned, err := service.IsGuestBanned(uuid)
	if err != nil {
		HandleError(c, err)
		return
	}

	if !isBanned {
		HandleError(c, errors.New("guest user is not banned"))
		return
	}

	// Unban the guest user
	if err := service.UnbanGuestUser(uuid); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Guest user unbanned successfully",
	})
}

// GetGuestUsers godoc
// @Summary Получение всех гостевых пользователей (только для администраторов)
// @Tags admin
// @Produce json
// @Success 200 {object} []models.GuestUser
// @Failure 500 {object} map[string]string
// @Router /admin/guests [get]
func GetGuestUsers(c *gin.Context) {
	guestUsers, err := service.GetGuestUsers()
	if err != nil {
		HandleError(c, err)
		return
	}
	if guestUsers == nil {
		HandleError(c, errs.ErrNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"guestUsers": guestUsers,
	})		
}

// GetGuestUser godoc
// @Summary Получение гостевого пользователя по UUID (только для администраторов)
// @Tags admin
// @Produce json
// @Param uuid path string true "Guest UUID"
// @Success 200 {object} models.GuestUser
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/guests/{uuid} [get]
func GetGuestUser(c *gin.Context) {
	uuid := c.Param("uuid")

	// Get the guest user
	guestUser, err := service.GetGuestUser(uuid)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"guestUser": guestUser,
	})	
}

// UpdateReport godoc
// @Summary Обновление жалобы (только для администраторов)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "Report ID"
// @Param report body models.UpdateReport true "Report object"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/reports/{id} [put]
func UpdateReport(c *gin.Context) {
	// Get user role from context
	userRole, exists := c.Get(middleware.RoleCtx)
	if !exists || userRole != "admin" {
		HandleError(c, errs.ErrForbidden)
		return
	}

	reportID, err := strconv.Atoi(c.Param("id"))
	if err != nil || reportID <= 0 {
		HandleError(c, errs.ErrInvalidId)
		return
	}

	// First check if report exists
	_, err = service.GetReport(reportID)
	if err != nil {
		HandleError(c, err)
		return
	}

	var updateReq models.UpdateReport
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		HandleError(c, err)
		return
	}

	// Update the report
	if err := service.UpdateReport(reportID, updateReq); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Report updated successfully",
	})
}

// GetReport godoc
// @Summary Получение жалобы по ID (только для администраторов)
// @Tags admin
// @Produce json
// @Param id path int true "Report ID"
// @Success 200 {object} models.Report
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/reports/{id} [get]
func GetReport(c *gin.Context) {
	reportID, err := strconv.Atoi(c.Param("id"))
	if err != nil || reportID <= 0 {
		HandleError(c, errs.ErrInvalidId)
		return
	}

	// First check if report exists
	report, err := service.GetReport(reportID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"report": report,
	})	
}