package service

import (
	"context"
	"fmt"
	"restaurant-service/models"
	"restaurant-service/utils"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MenuService interface {
	AddDishes(context.Context, *models.AddDishRequest) (*models.AddDishesResponse, error)
	IsRestaurantExist(context.Context, string) (bool, error)
}

type MenuServiceOptions struct {
	RestaurantCollection *mongo.Collection
	MenuCollection       *mongo.Collection
	Ctx                  context.Context
}

type MenuServiceImpl struct {
	opts MenuServiceOptions
}

func NewMenuService(opts MenuServiceOptions) MenuService {
	return &MenuServiceImpl{
		opts: opts,
	}
}

func (ms *MenuServiceImpl) AddDishes(ctx context.Context, req *models.AddDishRequest) (*models.AddDishesResponse, error) {
	isRestaurantExist, err := ms.IsRestaurantExist(ctx, req.RestaurantId)
	if err != nil {
		return &models.AddDishesResponse{}, err
	}

	if !isRestaurantExist {
		return &models.AddDishesResponse{}, fmt.Errorf("restaurant doen't exists")
	}

	dishesNames := []string{}
	insertOps := []mongo.WriteModel{} 

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, dish := range req.Dishes {

		wg.Add(1)
		d := dish
		go func() {

			defer wg.Done()
			imageUrl, err := utils.UploadImageToCloudinary(ctx, dish.Image)
			if err != nil {
				fmt.Printf("Error uploading image for dish %s", dish.Name)
				return
			}


			newDish := models.AddDishesInsert{
				Name:        d.Name,
				Price:       d.Price,
				Category:    d.Category,
				IsAvailable: d.IsAvailable,
				Image:       imageUrl,
				Rating:      d.Rating,
				Serves:      d.Serves,
				Discount:    d.Discount,
				CreatedAt:   time.Now(),
				UpdateAt:    time.Now(),
			}

			mu.Lock()
			insertOp := mongo.NewInsertOneModel().SetDocument(newDish)
			insertOps = append(insertOps, insertOp)
			dishesNames = append(dishesNames, dish.Name)
			mu.Unlock()
		}()
	}
	wg.Wait()

	if len(insertOps) > 0 {
		_, err = ms.opts.MenuCollection.BulkWrite(ctx, insertOps)
		if err != nil {
			return &models.AddDishesResponse{}, fmt.Errorf("error inserting dishes: %v", err)
		}
	}

	return &models.AddDishesResponse{
		RestaurantId: req.RestaurantId,
		DishesAdded:  dishesNames,
	}, nil

}

func (ms *MenuServiceImpl) IsRestaurantExist(ctx context.Context, restuarantId string) (bool, error) {

	if restuarantId == "" {
		return false, fmt.Errorf("restaurant Id is empty")
	}

	restuarantIdPrimitve, err := primitive.ObjectIDFromHex(restuarantId)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": restuarantIdPrimitve}

	var restaurant bson.M
	err = ms.opts.RestaurantCollection.FindOne(ctx, filter).Decode(&restaurant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, fmt.Errorf("no document exists")
		}
		return false, err
	}
	return true, nil
}
