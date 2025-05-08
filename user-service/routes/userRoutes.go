package routes

import (
	"user-service/controller"
	"user-service/service"

	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	UserController *controller.UserController
}

func NewUserRouteController(userController *controller.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) RegisterUser(rg *gin.RouterGroup, userService service.UserService) {
	router := rg.Group("users")
	router.POST("/registerUser", uc.UserController.RegisterUser)
	// router.Use(middleware.JWTAuthMiddleware())
	router.POST("/login", uc.UserController.LoginUser)
}
