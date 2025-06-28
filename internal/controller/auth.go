package controller

import (
	"github.com/hadisjane/confessly/internal/models"
	"github.com/hadisjane/confessly/internal/service"
	"github.com/hadisjane/confessly/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Регистрация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserRegister true "User object"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var u models.UserRegister

	if err := c.ShouldBindJSON(&u); err != nil {
		HandleError(c, err)
		return
	}

	if err := service.CreateUser(u); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})

}

// Login godoc
// @Summary Авторизация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserLogin true "User object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var u models.UserLogin

	if err := c.ShouldBindJSON(&u); err != nil {
		HandleError(c, err)
		return
	}

	user, err := service.GetUserByUsernameAndPassword(u.Username, u.Password)
	if err != nil {
		HandleError(c, err)
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}
