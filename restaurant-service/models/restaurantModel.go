package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Add restuarants
type AddRestaurantRequest struct {
	Name        string  `json:"name" bson:"name" binding:"required"`
	Address     string  `json:"address" bson:"address" binding:"required"`
	IsOpen      bool    `json:"isOpen" bson:"isOpen"`
	OpeningTime string  `json:"openingTime" bson:"openingTime" binding:"required"`
	ClosingTime string  `json:"closingTime" bson:"closingTime" binding:"required"`
	Location    GeoJSON `json:"location" bson:"location"`
}

type RestuarantInsertModel struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name" binding:"required"`
	Address     string             `json:"address" bson:"address" binding:"required"`
	IsOpen      bool               `json:"isOpen" bson:"isOpen"`
	OpeningTime time.Time          `json:"openingTime" bson:"openingTime" binding:"required"`
	ClosingTime time.Time          `json:"closingTime" bson:"closingTime" binding:"required"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdateAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
	Location    GeoJSON            `json:"location" bson:"location"`
}

type AddRestaurantResponse struct {
	Id      any    `json:"restaurantId" bson:"_id"`
	Message string `json:"message" bson:"message"`
}

// Get Restaurants
type GetRestauranstByLocationRequest struct {
	UserLocation GeoJSON `json:"userLocation" bson:"userLocation" binding:"required"`
	SortBy       string  `json:"sortBy" bson:"sortBy"`
	OrderBy      int     `json:"orderBy" bson:"orderBy"`
	Distance     float64 `json:"distance" bson:"distance"`
	IsOpen       bool    `json:"isOpen" bson:"isOpen"`
}

type GetRestauranstByLocation struct {
	RestaurantId   primitive.ObjectID `json:"_id" bson:"_id"`
	RestaurantName string             `json:"name" bson:"name"`
	Address        string             `json:"address" bson:"address"`
	OpeningTime    any             `json:"openingTime" bson:"openingTime"`
	ClosingTime    any             `json:"closingTime" bson:"closingTime"`
	IsOpen         bool               `json:"isOpen" bson:"isOpen"`
	DistanceInKms  float64            `json:"distanceInKms" bson:"distanceInKms"`
}

type GetRestauranstByLocationResponse struct {
	Restaurants []GetRestauranstByLocation `json:"restaurants" bson:"restaurants"`
}

// menuItems
type MenuItem struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	RestaurantId primitive.ObjectID `json:"restaurant_id" bson:"restaurant_id" bindind:"required"`
	Name         string             `json:"name" bson:"name" bindind:"required"`
	Description  string             `json:"description" bson:"description"`
	Category     string             `json:"category" bson:"category"`
	Price        float64            `json:"price" bson:"price" binding:"required"`
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
