package controller

import (
	"net/http"
	"restaurant-service/models"
	"restaurant-service/service"
	"restaurant-service/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuController struct {
	MenuService service.MenuService
}

func NewMenuController(menuService service.MenuService) *MenuController {
	return &MenuController{menuService}
}

func (mc *MenuController) AddDishes(c *gin.Context) {
	ctx := c.Request.Context()

	form, err := c.MultipartForm()
	if err != nil {
		utils.ResponseHandler(c, false, http.StatusBadRequest, "Invalid multipart form", err)
	}
	restaurantId := form.Value["restaurant_id"]

	names := form.Value["name"]
	prices := form.Value["price"]
	categories := form.Value["category"]
	isAvailables := form.Value["isAvailable"]
	ratings := form.Value["rating"]
	servesList := form.Value["serves"]
	discounts := form.Value["discount"]
	images := form.File["image"]

	if len(names) != len(prices) || len(names) != len(images) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Inconsistent number of dish fields"})
		return
	}

	var req models.AddDishRequest
	if len(restaurantId) > 0 {
		req.RestaurantId = restaurantId[0]
	}

	for i := range names {
		price, _ := strconv.ParseFloat(prices[i], 64)
		available, _ := strconv.ParseBool(isAvailables[i])
		rating, _ := strconv.ParseFloat(ratings[i], 64)
		serves, _ := strconv.Atoi(servesList[i])
		discount, _ := strconv.ParseFloat(discounts[i], 64)

		dish := models.AddDishes{
			Name:        names[i],
			Price:       price,
			Category:    categories[i],
			IsAvailable: available,
			Image:       images[i],
			Rating:      rating,
			Serves:      serves,
			Discount:    discount,
		}
		req.Dishes = append(req.Dishes, dish)
	}

	response, err := mc.MenuService.AddDishes(ctx, &req)
	if err != nil {
		utils.ResponseHandler(c, false, http.StatusBadRequest, "Unable to add dishes", err)
	}

	utils.ResponseHandler(c, true, http.StatusOK, "Successfully added dishes", response)
}
