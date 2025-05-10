package routes

import (
	"restaurant-service/controller"
	"restaurant-service/service"

	"github.com/gin-gonic/gin"
)

type NewRestaurantRoutesController struct {
	RestaurantController *controller.RestaurantController
}

func NewRestuarantRouteController(restuarantController *controller.RestaurantController) NewRestaurantRoutesController {
	return NewRestaurantRoutesController{
		RestaurantController: restuarantController,
	}
}

func (rc *NewRestaurantRoutesController) RestaurantRoutes(rg *gin.RouterGroup, restaurantService service.RestaurantService) {
	router := rg.Group("restaurants")
	router.POST("/restaurants", rc.RestaurantController.AddRestaurant)
	router.GET("/restaurants", rc.RestaurantController.GetRestaurantsByLocation)
}
