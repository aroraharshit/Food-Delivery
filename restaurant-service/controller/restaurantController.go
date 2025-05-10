package controller

import (
	"net/http"
	"restaurant-service/models"
	"restaurant-service/service"

	"github.com/gin-gonic/gin"
)

type RestaurantController struct {
	RestaurantService service.RestaurantService
}

func NewRestaurantController(restaurantService service.RestaurantService) *RestaurantController {
	return &RestaurantController{
		RestaurantService: restaurantService,
	}
}

func (rs *RestaurantController) AddRestaurant(c *gin.Context) {
	var req models.AddRestaurantRequest
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request" + err.Error()})
		return
	}

	response, err := rs.RestaurantService.AddRestaurant(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (rs *RestaurantController) GetRestaurantsByLocation(c *gin.Context) {
	var req models.GetRestauranstByLocationRequest
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request" + err.Error()})
	}

	response, err := rs.RestaurantService.GetRestaurantByLocation(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
