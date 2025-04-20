package repositories

import (
	"sample-web/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	CreateUser(user models.User) (models.User, error)
	FindUserByEmail(email string) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
}

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (userRepository *userRepository) CreateUser(user models.User) (models.User, error) {
	return models.User{}, nil
}

func (userRepository *userRepository) FindUserByEmail(email string) (models.User, error) {
	return models.User{}, nil
}

func (userRepository *userRepository) UpdateUser(user models.User) (models.User, error) {
	return models.User{}, nil
}
