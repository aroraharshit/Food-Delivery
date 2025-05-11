package controller

import (
	"net/http"
	"restaurant-service/models"
	"restaurant-service/service"
	"restaurant-service/utils"

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
		utils.ResponseHandler(c, false, http.StatusBadRequest, "Invalid request", err)
		return
	}

	response, err := rs.RestaurantService.AddRestaurant(ctx, &req)
	if err != nil {
		utils.ResponseHandler(c, false, http.StatusBadRequest, "Failed to insert the restaurant", err)
		return
	}

	utils.ResponseHandler(c, true, http.StatusOK, "Added restuarant successfully", response)
}

// GetRestaurants godoc
// @Summary      Get list of restuarants
// @Description  Returns all restaurants
// @Tags         restuarants
// @Produce      json
// @Param        request  body    models.GetRestauranstByLocationRequest  true  "Location and Filter Parameters"
// @Success      200  {array}  models.GetRestauranstByLocationResponse
// @Router       /v1/restaurants/getRestaurants [post]
func (rs *RestaurantController) GetRestaurantsByLocation(c *gin.Context) {
	var req models.GetRestauranstByLocationRequest
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseHandler(c, false, http.StatusBadRequest, "Invalid Request", err)
		return
	}

	response, err := rs.RestaurantService.GetRestaurantByLocation(ctx, &req)
	if err != nil {
		utils.ResponseHandler(c, false, http.StatusBadRequest, "Failed to fetch restuarant by location", err)
		return
	}

	utils.ResponseHandler(c, true, http.StatusOK, "Restuarants fetched successfully", response)
}
