package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"restaurant-service/models"
	"restaurant-service/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RestaurantSevice interface {
	AddRestaurant(context.Context, *models.AddRestaurantRequest) (*models.AddRestaurantResponse, error)
	GetRestaurants(context.Context, *models.GetRestauranstByLocationRequest) (*models.GetRestauranstByLocationResponse, error)
}

type RestaurantServiceOptions struct {
	ctx                  context.Context
	RestaurantCollection *mongo.Collection
}

type RestaurantService struct {
	opts RestaurantServiceOptions
}

func NewRestaurantService(opts RestaurantServiceOptions) RestaurantService {
	return RestaurantService{
		opts: opts,
	}
}

func (rs *RestaurantService) AddRestaurant(ctx context.Context, req *models.AddRestaurantRequest) (*models.AddRestaurantResponse, error) {

	openingTime, err := utils.StringTimeToISO(req.OpeningTime)
	if err != nil {
		return &models.AddRestaurantResponse{}, err
	}

	closingTime, err := utils.StringTimeToISO(req.ClosingTime)
	if err != nil {
		return &models.AddRestaurantResponse{}, err
	}

	insertDocument := models.RestuarantInsertModel{
		Id:          primitive.NewObjectID(),
		Name:        req.Name,
		Address:     req.Address,
		IsOpen:      req.IsOpen,
		OpeningTime: openingTime,
		ClosingTime: closingTime,
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
		Location:    req.Location,
	}

	res, err := rs.opts.RestaurantCollection.InsertOne(ctx, insertDocument)
	if err != nil {
		return &models.AddRestaurantResponse{}, err
	}
	return &models.AddRestaurantResponse{Message: "Restuarant added", Id: res.InsertedID}, nil
}

func (rs *RestaurantService) GetRestaurantByLocation(ctx context.Context, req *models.GetRestauranstByLocationRequest) (*models.GetRestauranstByLocationResponse, error) {
	var response models.GetRestauranstByLocationResponse

	locationType := req.UserLocation.Type
	coordinates := req.UserLocation.Coordinates
	lat := coordinates[0]
	long := coordinates[1]
	distance := 5000

	pipeline := bson.A{bson.M{
		"$geoNear": bson.M{
			"near": bson.M{
				"type":        locationType,
				"coordinates": bson.A{lat, long},
			},
			"distanceField": "distance",
			"maxDistance":   distance,
			"spherical":     true,
		},
	},
		bson.M{
			"$addFields": bson.M{
				"distanceInKms": bson.M{
					"$round": bson.A{
						bson.M{
							"$divide": bson.A{"$distance", 1000},
						},
						1,
					},
				},
			},
		}}

	qry, _ := json.Marshal(pipeline)
	fmt.Println(string(qry))

	cursor, err := rs.opts.RestaurantCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
		return &models.GetRestauranstByLocationResponse{}, err
	}

	defer cursor.Close(ctx)

	var restaurants []models.GetRestauranstByLocation

	for cursor.Next(ctx) {
		var restaurant models.GetRestauranstByLocation
		if err := cursor.Decode(&restaurant); err != nil {
			log.Println("Error decoding cursor", err)
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	response = models.GetRestauranstByLocationResponse{
		Restaurants: restaurants,
	}
	return &response, nil

}
