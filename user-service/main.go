package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"user-service/controller"
	"user-service/routes"
	"user-service/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	if DBName == "" {
		return nil, fmt.Errorf("MONGO_DB_NAME not set")
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := connectMongoDB(ctx)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	userCollection := client.Database("DBName").Collection("users")

	userService := service.NewUserService(service.UserServiceOptions{
		Ctx:            ctx,
		UserCollection: userCollection,
	})

	userController := controller.NewUserController(userService)

	userRouteController := routes.NewUserRouteController(userController)

	router := gin.Default()
	api := router.Group("/v1")
	userRouteController.RegisterUser(api, userService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Server is running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
