package models

import (
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddDishRequest struct {
	RestaurantId string      `form:"restaurant_id"`
	Dishes       []AddDishes `form:"dishes[]"`
}

type AddDishes struct {
	RestaurantId string                `form:"restaurant_id"`
	Name         string                `form:"name"`
	Price        float64               `form:"price"`
	Category     string                `form:"category"`
	IsAvailable  bool                  `form:"isAvailable"`
	Image        *multipart.FileHeader `form:"image"`
	Rating       float64               `form:"rating"`
	Serves       int                   `form:"serves"`
	Discount     float64               `form:"discount"`
}

type AddDishesInsert struct {
	RestaurantId primitive.ObjectID `bson:"restaurant_id"`
	Name         string             `bson:"name"`
	Price        float64            `bson:"price"`
	Category     string             `bson:"category"`
	IsAvailable  bool               `bson:"isAvailable"`
	Image        string             `bson:"image"`
	Rating       float64            `bson:"rating"`
	Serves       int                `bson:"serves"`
	Discount     float64            `bson:"discount"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdateAt     time.Time          `bson:"updatedAt"`
}

type AddDishesInsertPayload struct {
	RestaurantId string            `bson:"restaurant_id"`
	Dishes       []AddDishesInsert `bson:"dishes"`
}

type AddDishesResponse struct {
	RestaurantId string   `json:"restaurant_id" bson:"restaurant_id"`
	DishesAdded  []string `json:"dishesAdded" bson:"dishesAdded"`
}
