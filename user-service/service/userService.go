package service

import (
	"context"
	"fmt"
	"log"
	"time"
	"user-service/model"
	"user-service/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(*model.RegisterUserRequest) (*model.RegisterUserResponse, error)
	UserExists(email string, mobile string) (bool, error)
	LoginUser(*model.LoginUserRequest) (*model.LoginUserResponse, error)
}

type UserServiceOptions struct {
	UserCollection *mongo.Collection
}

type userServiceImpl struct {
	opts UserServiceOptions
}

func NewUserService(opts UserServiceOptions) UserService {
	return &userServiceImpl{
		opts: opts,
	}
}

func (us *userServiceImpl) UserExists(email string, mobileNumber string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"$or": bson.A{
					bson.M{
						"mobile": mobileNumber,
					},
					bson.M{
						"email": email,
					},
				},
			},
		},
		bson.M{
			"$count": "count",
		},
		bson.M{
			"$addFields": bson.M{
				"userExists": bson.M{
					"$gt": bson.A{"$count", 0},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"userExists": 1,
				"_id":        0,
			},
		},
	}

	cursor, err := us.opts.UserCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	defer cursor.Close(ctx)

	var result []model.UserExists
	if err := cursor.All(ctx, &result); err != nil {
		log.Fatal(err)
		return false, err
	}

	if len(result) == 0 {
		return false, nil
	}

	return result[0].UserExists, nil
}

func (us *userServiceImpl) RegisterUser(req *model.RegisterUserRequest) (*model.RegisterUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userExist, err := us.UserExists(req.Email, req.Mobile)
	if err != nil {
		return nil, fmt.Errorf("Error checking if user exists or not: %v", err)
	}
	if userExist {
		return &model.RegisterUserResponse{
			Message: "User already exists",
			UserId:  "",
		}, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Error while hashing password %v", err)
	}

	user := model.RegistertUserInsertion{
		Email:     req.Email,
		Name:      req.Name,
		Mobile:    req.Mobile,
		Password:  string(hashedPassword),
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		Id:        primitive.NewObjectID(),
	}

	insertResult, err := us.opts.UserCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error inserting user into database %v", err)
	}

	response := &model.RegisterUserResponse{
		Message: "User registered successful",
		UserId:  insertResult.InsertedID.(primitive.ObjectID).Hex(),
	}

	return response, nil
}

func (us *userServiceImpl) GetPasswordByNumber(mobile string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result struct {
		Password string `bson:"password"`
	}

	err := us.opts.UserCollection.FindOne(ctx, bson.M{"mobile": mobile}, options.FindOne().SetProjection(bson.M{"password": 1})).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("user with mobile %s not found", mobile)
		}
		return "", fmt.Errorf("error finding password of user by mobile: %v", err)
	}

	return result.Password, nil
}

func (us *userServiceImpl) LoginUser(req *model.LoginUserRequest) (*model.LoginUserResponse, error) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userExist, err := us.UserExists(req.Email, req.Mobile)
	if err != nil {
		return nil, fmt.Errorf("error checking if user exists or not: %v", err)
	}

	if !userExist {
		return &model.LoginUserResponse{
			Message: "User doesnt exists",
			Token:   "",
		}, nil
	}

	storedPassword, err := us.GetPasswordByNumber(req.Mobile)
	if err != nil {
		return nil, fmt.Errorf("error getting password from database %v", err)
	}

	correctPassword := utils.VerifyPassword(req.Password, storedPassword)
	if !correctPassword {
		return &model.LoginUserResponse{
			Message: "Incorrect Password",
			Token:   "",
		}, nil
	}

	token, err := utils.CreateToken(req.Mobile)
	if err != nil {
		return nil, fmt.Errorf("error generating token %v", err)
	}

	return &model.LoginUserResponse{
		Message: "Successfully Login",
		Token:   token,
	}, nil

}
