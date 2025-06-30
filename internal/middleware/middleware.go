package middleware

import (
	"github.com/hadisjane/confessly/internal/db"
	"github.com/hadisjane/confessly/internal/models"
	"github.com/hadisjane/confessly/internal/service"
	"github.com/hadisjane/confessly/logger"
	"github.com/hadisjane/confessly/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	UserIDCtx           = "userID"
	UsernameCtx         = "username"
	RoleCtx             = "role"
	GuestUUIDCtx        = "guestUUID"
)

func CheckUserAuthentication(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "empty auth header",
		})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid auth header",
		})
		return
	}

	if len(headerParts[1]) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "token is empty",
		})
		return
	}

	accessToken := headerParts[1]

	claims, err := utils.ParseToken(accessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get database connection
	db := db.GetDB()
	if db == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	// Check if user is banned
	isBanned, err := service.IsUserBanned(claims.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to check user status"})
		return
	}

	if isBanned {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "your account has been banned",
		})
		return
	}

	c.Set(UserIDCtx, claims.UserID)
	c.Set(UsernameCtx, claims.Username)
	c.Set(RoleCtx, claims.Role)
	c.Next()
}

func CheckAdminAuthentication(c *gin.Context) {
	role, exists := c.Get(RoleCtx)
	logger.Info.Printf("Role from context: %v, exists: %v", role, exists)
	
	roleStr, ok := role.(string)
	if !ok || roleStr != "admin" {
		 c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		  	  "error": "forbidden: admin role required",
		 })
		 return
	}
	c.Next()
}

func TryParseUserContext(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.Next()
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" || len(headerParts[1]) == 0 {
		c.Next()
		return
	}

	accessToken := headerParts[1]

	claims, err := utils.ParseToken(accessToken)
	if err != nil {
		c.Next()
		return
	}

	// Check if user is banned
	isBanned, err := service.IsUserBanned(claims.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to check user status"})
		return
	}

	if isBanned {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "your account has been banned",
		})
		return
	}

	c.Set(UserIDCtx, claims.UserID)
	c.Set(UsernameCtx, claims.Username)
	c.Set(RoleCtx, claims.Role)

	c.Next()
}

func GuestUUIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip guest check if user is already authenticated
		if userID, exists := c.Get(UserIDCtx); exists && userID != nil {
			c.Next()
			return
		}

		// Check for guest UUID in cookie
		guestUUID, err := c.Cookie("guest_uuid")
		if err != nil || guestUUID == "" {
			// If no cookie, generate a new guest UUID and create user
			createNewGuestUser(c)
			return
		}

		// Check if guest exists in database
		_, err = service.GetGuestUser(guestUUID)
		if err != nil {
			// If guest doesn't exist, create a new one
			logger.Info.Printf("Guest user %s not found, creating new one", guestUUID)
			createNewGuestUser(c)
			return
		}

		// Check if guest is banned
		isBanned, err := service.IsGuestBanned(guestUUID)
		if err != nil {
			logger.Error.Printf("Failed to check guest ban status: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to check guest status"})
			return
		}

		if isBanned {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "your guest account has been banned",
			})
			return
		}

		c.Set(GuestUUIDCtx, guestUUID)
		c.Next()
	}
}

func createNewGuestUser(c *gin.Context) {
	// Generate a new guest UUID
	guestUUID := uuid.New().String()
	
	// Set cookie with new UUID
	c.SetCookie("guest_uuid", guestUUID, 30*24*60*60, "/", "", false, true)

	// Create and save new guest user
	guest := models.GuestUser{
		UUID:   guestUUID,
		Banned: false,
	}

	if err := service.CreateGuestUser(guest); err != nil {
		logger.Error.Printf("Failed to create guest user: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create guest user"})
		return
	}

	c.Set(GuestUUIDCtx, guestUUID)
}
