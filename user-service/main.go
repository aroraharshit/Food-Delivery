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

func ConnectToMongoDB() *mongo.Collection {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURL := os.Getenv("MONGODB_URL")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	duration := time.Since(startTime)

	fmt.Printf("Connected to MongoDB in %s\n", duration)

	dbName := os.Getenv("MONGO_DB_NAME")
	collection := client.Database(dbName).Collection("users")
	fmt.Println("Connected to MongoDB")
	return collection
}

func main() {
	router := gin.Default()

	userCollection := ConnectToMongoDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// user collection
	userService := service.NewUserService(service.UserServiceOptions{
		UserCollection: userCollection,
	})
	userController := controller.NewUserController(userService)
	userRouterController := routes.NewUserRouteController(userController)

	api := router.Group("/v1")
	userRouterController.RegisterUser(api, userService)

	fmt.Println("Server started at :", port)
	router.Run(":" + port)

}
