package routes

import (
	"restaurant-service/controller"
	"restaurant-service/service"

	"github.com/gin-gonic/gin"
)

type MenuRoutesController struct {
	MenuController *controller.MenuController
}

func NewMenuRouteController(menuController *controller.MenuController) MenuRoutesController{
	return MenuRoutesController{MenuController:menuController }
}

func (mc *MenuRoutesController) MenuRoutes(rg *gin.RouterGroup,menuService service.MenuService){
	router:=rg.Group("menu")
	router.POST("/addDishes",mc.MenuController.AddDishes)
}