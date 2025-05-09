package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddRestaurantRequest struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Address     string             `json:"address" bson:"address"`
	IsOpen      bool               `json:"isOpen" bson:"isOpen"`
	OpeningTime time.Time          `json:"openingTime" bson:"openingTime"`
	ClosingTime time.Time          `json:"closingTime" bson:"closingTime"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdateAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
	Location    GeoJSON            `json:"location" bson:"location"`
}

type AddRestaurantResponse struct {
	Id      primitive.ObjectID `json:"restaurantId" bson:"_id"`
	Message string             `json:"message" bson:"message"`
}

type MenuItem struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	RestaurantId primitive.ObjectID `json:"restaurant_id" bson:"restaurant_id"`
	Name         string             `json:"name" bson:"name"`
	Description  string             `json:"description" bson:"description"`
	Category     string             `json:"category" bson:"category"`
	Price        float64            `json:"price" bson:"price"`
	IsAvailable  bool               `json:"isAvailable" bson:"isAvailable"`
	IsVeg        bool               `json:"isVeg" bson:"isVeg"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updateAt" bson:"updatedAt"`
	Tags         []string           `json:"tags" bson:"tags"`
}

type GeoJSON struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}
