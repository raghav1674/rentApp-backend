package repositories

import (
	"sample-web/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	CreateUser(ctx *gin.Context, user models.User) (models.User, error)
	FindUserByEmail(ctx *gin.Context, email string) (models.User, error)
	UpdateUser(ctx *gin.Context, user models.User) (models.User, error)
}

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (userRepository *userRepository) CreateUser(ctx *gin.Context, user models.User) (models.User, error) {
	usersCollection := userRepository.db.Collection("users")
	result, err := usersCollection.InsertOne(ctx, user)
	if err != nil {
		return models.User{}, err
	}
	var createdUser models.User

	err = usersCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdUser)
	if err != nil {
		return models.User{}, err
	}
	return createdUser, nil
}

func (userRepository *userRepository) FindUserByEmail(ctx *gin.Context, email string) (models.User, error) {
	usersCollection := userRepository.db.Collection("users")
	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (userRepository *userRepository) UpdateUser(ctx *gin.Context, user models.User) (models.User, error) {
	usersCollection := userRepository.db.Collection("users")
	_, err := usersCollection.UpdateOne(ctx, bson.M{"_id": user.Id}, bson.M{"$set": user})
	if err != nil {
		return models.User{}, err
	}
	var updatedUser models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": user.Id}).Decode(&updatedUser)
	if err != nil {
		return models.User{}, err
	}
	return updatedUser, nil
}
