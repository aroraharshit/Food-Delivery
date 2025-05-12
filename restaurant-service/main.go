package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"restaurant-service/controller"
	_ "restaurant-service/docs"
	"restaurant-service/routes"
	"restaurant-service/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	utils "restaurant-service/utils"
)

var (
	DBName string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func connectMongoDB(ctx context.Context) (*mongo.Client, error) {
	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		return nil, fmt.Errorf("MONGODB_URL not set")
	}

	clientOpts := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %w", err)
	}

	DBName = os.Getenv("MONGO_DB_NAME")
	fmt.Println("DBname", DBName)
	if DBName == "" {
		return nil, fmt.Errorf("MONGO_DB_NAME not set")
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

// @title           Food Delivery API
// @version         1.0
// @description     API Server for Food Delivery App
// @host            localhost:8080
// @BasePath        /

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	client, err := connectMongoDB(ctx)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	utils.IntiCloudinary()

	//Collection Initialisation
	restaurantCollection := client.Database(DBName).Collection("restaurant")
	menuCollection := client.Database(DBName).Collection("menu")

	//services initialisation
	restaurantService := service.NewRestaurantService(
		service.RestaurantServiceOptions{
			Ctx:                  ctx,
			RestaurantCollection: restaurantCollection,
		},
	)
	menuService := service.NewMenuService(service.MenuServiceOptions{
		RestaurantCollection: restaurantCollection,
		MenuCollection:       menuCollection,
		Ctx:                  ctx,
	},
	)

	//Controller initialisation
	restaurantController := controller.NewRestaurantController(restaurantService)
	menuContoller := controller.NewMenuController(menuService)

	//Routes initialisation
	restaurantRouteController := routes.NewRestuarantRouteController(restaurantController)
	menuRouteController := routes.NewMenuRouteController(menuContoller)

	router := gin.Default()
	
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/v1")
	restaurantRouteController.RestaurantRoutes(api, restaurantService)
	menuRouteController.MenuRoutes(api, menuService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Server is running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
