package repositories

import (
	"context"
	"sample-web/models"
	"sample-web/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	FindUserByEmail(ctx context.Context, email string) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
}

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (userRepository *userRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {

	ctx, span := utils.Tracer().Start(ctx, "UserRepository.CreateUser")
	defer span.End()

	span.AddEvent("mongo.InsertOne", trace.WithAttributes(
		attribute.String("collection", "users"),
		attribute.String("operation", "insert_one"),
		attribute.String("email", user.Email),
	))

	usersCollection := userRepository.db.Collection("users")
	_, err := usersCollection.InsertOne(ctx, user)
	if err != nil {
		span.RecordError(err)
		span.AddEvent("UserCreationFailed")
		return models.User{}, err
	}

	span.AddEvent("UserCreated")

	return userRepository.FindUserByEmail(ctx, user.Email)
}

func (userRepository *userRepository) FindUserByEmail(ctx context.Context, email string) (models.User, error) {

	ctx, span := utils.Tracer().Start(ctx, "UserRepository.FindUserByEmail")
	defer span.End()

	usersCollection := userRepository.db.Collection("users")
	var user models.User

	span.AddEvent("mongo.FindOne", trace.WithAttributes(
		attribute.String("collection", "users"),
		attribute.String("operation", "find_one"),
		attribute.String("email", email),
	))

	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		span.RecordError(err)
		return models.User{}, err
	}

	span.AddEvent("UserFound")

	return user, nil
}

func (userRepository *userRepository) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
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
