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
)

type UserService interface {
	RegisterUser(context.Context,*model.RegisterUserRequest) (*model.RegisterUserResponse, error)
	UserExists(ctx context.Context,email string, mobile string) (bool, error)
	LoginUser(context.Context,*model.LoginUserRequest) (*model.LoginUserResponse, error)
}

type UserServiceOptions struct {
	ctx context.Context
	UserCollection *mongo.Collection
}

type UserServiceImpl struct {
	opts UserServiceOptions
}

func NewUserService(opts UserServiceOptions) UserService {
	return &UserServiceImpl{
		opts: opts,
	}
}

func (us *UserServiceImpl) UserExists(ctx context.Context,email string, mobileNumber string) (bool, error) {
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

func (us *UserServiceImpl) RegisterUser(ctx context.Context,req *model.RegisterUserRequest) (*model.RegisterUserResponse, error) {
	userExist, err := us.UserExists(ctx,req.Email, req.Mobile)
	if err != nil {
		return nil, fmt.Errorf("Error checking if user exists or not: %v", err)
	}
	if userExist {
		return &model.RegisterUserResponse{
			Message: "User already exists",
			UserId:  "",
		}, nil
	}

	hashedPassword, err := utils.HashedPassword(req.Password)
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

func (us *UserServiceImpl) GetPasswordByNumber(ctx context.Context,mobile string) (string, error) {
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

func (us *UserServiceImpl) LoginUser(ctx context.Context,req *model.LoginUserRequest) (*model.LoginUserResponse, error) {
	userExist, err := us.UserExists(ctx,req.Email, req.Mobile)
	if err != nil {
		return nil, fmt.Errorf("error checking if user exists or not: %v", err)
	}

	if !userExist {
		return &model.LoginUserResponse{
			Message: "User doesnt exists",
			Token:   "",
		}, nil
	}

	storedPassword, err := us.GetPasswordByNumber(ctx,req.Mobile)
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
