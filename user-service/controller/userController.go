package controller

import (
	"net/http"
	"user-service/model"
	"user-service/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) RegisterUser(c *gin.Context) {
	var req model.RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request:" + err.Error()})
		return
	}

	response, err := uc.userService.RegisterUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (uc *UserController) LoginUser(c *gin.Context) {
	var req model.LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request:" + err.Error()})
		return
	}

	response, err := uc.userService.LoginUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, response)
}
