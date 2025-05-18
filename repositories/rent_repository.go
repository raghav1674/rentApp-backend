package repositories

import (
	"context"
	"fmt"
	"sample-web/models"
	"sample-web/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type RentRepository interface {
	CreateRent(ctx context.Context, rent models.Rent) (models.Rent, error)
	FindRentById(ctx context.Context, userId string, rentId string) (models.Rent, error)
	GetAllRents(ctx context.Context, userId string, userRole models.UserRole) ([]models.Rent, error)
	UpdateRent(ctx context.Context, userId string, rent models.Rent) (models.Rent, error)
}

type rentRepository struct {
	db *mongo.Database
}

func NewRentRepository(db *mongo.Database) RentRepository {
	return &rentRepository{
		db: db,
	}
}

func (rentRepository *rentRepository) CreateRent(ctx context.Context, rent models.Rent) (models.Rent, error) {

	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx, "RentRepository.CreateRent")
	defer span.End()

	span.AddEvent("mongo.InsertOne", trace.WithAttributes(
		attribute.String("collection", "rents"),
		attribute.String("operation", "insert_one"),
	))

	rentsCollection := rentRepository.db.Collection("rents")
	result, err := rentsCollection.InsertOne(ctx, rent)
	if err != nil {
		span.RecordError(err)
		span.AddEvent("RentCreationFailed")
		return models.Rent{}, err
	}

	span.AddEvent("RentCreated")

	id := result.InsertedID.(bson.ObjectID).Hex()

	log.Info(spanCtx, fmt.Sprintf("Rent created with ID: %s", id))

	return rentRepository.FindRentById(ctx, rent.LandLord.Id.Hex(), id)
}

func (rentRepository *rentRepository) FindRentById(ctx context.Context, userId string, rentId string) (models.Rent, error) {

	log := utils.GetLogger()
	_, span := log.Tracer().Start(ctx, "RentRepository.FindRentById")
	defer span.End()

	rentsCollection := rentRepository.db.Collection("rents")
	var rent models.Rent

	span.AddEvent("mongo.FindOne", trace.WithAttributes(
		attribute.String("collection", "rents"),
		attribute.String("operation", "find_one"),
		attribute.String("_id", rentId),
	))

	objectID, err := bson.ObjectIDFromHex(rentId)
	if err != nil {
		span.RecordError(err)
		return models.Rent{}, err
	}

	userObjectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		span.RecordError(err)
		return models.Rent{}, err
	}

	query := bson.M{
		"_id": objectID,
		"$or": []bson.M{
			{"landlord._id": userObjectId},
			{"tenant._id": userObjectId},
		},
	}

	err = rentsCollection.FindOne(ctx, query).Decode(&rent)

	if err != nil {

		span.RecordError(err)
		return models.Rent{}, err
	}

	span.AddEvent("RentFound")

	return rent, nil
}

func (rentRepository *rentRepository) GetAllRents(ctx context.Context, userId string, userRole models.UserRole) ([]models.Rent, error) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx, "RentRepository.GetAllRents")
	defer span.End()

	rentsCollection := rentRepository.db.Collection("rents")

	span.AddEvent("mongo.Find", trace.WithAttributes(
		attribute.String("collection", "rents"),
		attribute.String("operation", "find"),
		attribute.String("user_id", userId),
	))

	userObjectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	var query interface{}

	if userRole == models.LandLord {
		query = bson.M{"landlord._id": userObjectId}
	} else {
		query = bson.M{"tenant._id": userObjectId}
	}

	log.Info(spanCtx, query)

	cursor, err := rentsCollection.Find(ctx, query)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var rents []models.Rent
	if err := cursor.All(ctx, &rents); err != nil {
		span.RecordError(err)
		return nil, err
	}

	log.Info(spanCtx, fmt.Sprintf("Found %d rents", len(rents)))

	span.AddEvent("RentsFound")
	return rents, nil
}

func (rentRepository *rentRepository) UpdateRent(ctx context.Context, userId string, rent models.Rent) (models.Rent, error) {

	_, span := utils.Tracer().Start(ctx, "RentRepository.UpdateRent")
	defer span.End()

	rentsCollection := rentRepository.db.Collection("rents")

	span.AddEvent("mongo.UpdateOne", trace.WithAttributes(
		attribute.String("collection", "rents"),
		attribute.String("operation", "find_one"),
		attribute.String("_id", rent.Id.Hex()),
	))

	rentObjectId, err := bson.ObjectIDFromHex(rent.Id.Hex())
	if err != nil {
		span.RecordError(err)
		return models.Rent{}, err
	}

	userObjectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		span.RecordError(err)
		return models.Rent{}, err
	}

	query := bson.M{
		"_id":          rentObjectId,
		"landlord._id": userObjectId,
	}

	_, err = rentsCollection.UpdateOne(ctx, query, bson.M{"$set": rent})
	if err != nil {
		span.RecordError(err)
		return models.Rent{}, err
	}
	return rentRepository.FindRentById(ctx, userId, rent.Id.Hex())
}
