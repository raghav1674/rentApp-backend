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
	FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (models.User, error)
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

	_, span := utils.Tracer().Start(ctx, "UserRepository.CreateUser")
	defer span.End()

	span.AddEvent("mongo.InsertOne", trace.WithAttributes(
		attribute.String("collection", "users"),
		attribute.String("operation", "insert_one"),
		attribute.String("phone_number", user.PhoneNumber),
	))

	usersCollection := userRepository.db.Collection("users")
	_, err := usersCollection.InsertOne(ctx, user)
	if err != nil {
		span.RecordError(err)
		span.AddEvent("UserCreationFailed")
		return models.User{}, err
	}

	span.AddEvent("UserCreated")

	return userRepository.FindUserByPhoneNumber(ctx, user.PhoneNumber)
}

func (userRepository *userRepository) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (models.User, error) {

	_, span := utils.Tracer().Start(ctx, "UserRepository.FindUserByPhoneNumber")
	defer span.End()

	usersCollection := userRepository.db.Collection("users")
	var user models.User

	span.AddEvent("mongo.FindOne", trace.WithAttributes(
		attribute.String("collection", "users"),
		attribute.String("operation", "find_one"),
		attribute.String("phone_number", phoneNumber),
	))

	err := usersCollection.FindOne(ctx, bson.M{"phone_number": phoneNumber}).Decode(&user)

	if err != nil {
		span.RecordError(err)
		return models.User{}, err
	}

	span.AddEvent("UserFound")

	return user, nil
}

func (userRepository *userRepository) UpdateUser(ctx context.Context, user models.User) (models.User, error) {

	_, span := utils.Tracer().Start(ctx, "UserRepository.UpdateUser")
	defer span.End()

	usersCollection := userRepository.db.Collection("users")

	span.AddEvent("mongo.UpdateOne", trace.WithAttributes(
		attribute.String("collection", "users"),
		attribute.String("operation", "find_one"),
		attribute.String("id", user.Id.Hex()),
		attribute.String("phone_number", user.PhoneNumber),
	))

	_, err := usersCollection.UpdateOne(ctx, bson.M{"_id": user.Id}, bson.M{"$set": user})
	if err != nil {
		span.RecordError(err)
		return models.User{}, err
	}
	return userRepository.FindUserByPhoneNumber(ctx, user.PhoneNumber)
}
