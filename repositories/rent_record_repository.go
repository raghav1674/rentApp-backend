package repositories

import (
	"context"
	"fmt"
	"sample-web/models"
	"sample-web/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RentRecordRepository interface {
	CreateRentRecord(ctx context.Context, rentRecord models.RentRecord) (models.RentRecord, error)
	GetRentRecordById(ctx context.Context, rentRecordId string) (models.RentRecord, error)
	GetAllRentRecords(ctx context.Context, userId string, userRole string, rentId string,) ([]models.RentRecord, error)
	UpdateRentRecord(ctx context.Context, userId string, rentRecordId string,rentRecord models.RentRecord) (models.RentRecord, error)
}

type rentRecordRepository struct {
	db *mongo.Database
}

func NewRentRecordRepository(db *mongo.Database) RentRecordRepository {
	return &rentRecordRepository{
		db: db,
	}
}

func (r *rentRecordRepository) CreateRentRecord(ctx context.Context, rentRecord models.RentRecord) (models.RentRecord, error) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx, "RentRecordRepository.CreateRentRecord")
	defer span.End()

	rentRecordcollection := r.db.Collection("rent_records")

	log.Info(spanCtx, "Inserting rent record into the database")

	result, err := rentRecordcollection.InsertOne(spanCtx, rentRecord)

	if err != nil {
		log.Error(spanCtx, "Error inserting rent record into the database")
		return models.RentRecord{}, err
	}
	log.Info(spanCtx, "Rent record inserted successfully")

	id := result.InsertedID.(bson.ObjectID).Hex()

	return r.GetRentRecordById(spanCtx, id)

}

func (r *rentRecordRepository) GetRentRecordById(ctx context.Context, rentRecordId string) (models.RentRecord, error) {

	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx, "RentRecordRepository.GetRentRecordById")
	defer span.End()

	rentRecordCollection := r.db.Collection("rent_records")
	log.Info(spanCtx, "Fetching rent record from the database")

	rentRecordObjectId, err := bson.ObjectIDFromHex(rentRecordId)
	if err != nil {
		log.Error(spanCtx, "Error converting rent record ID to ObjectID")
		return models.RentRecord{}, err
	}

	var rentRecord models.RentRecord
	err = rentRecordCollection.FindOne(spanCtx, bson.M{"_id": rentRecordObjectId}).Decode(&rentRecord)
	if err != nil {
		log.Error(spanCtx, "Error fetching rent record from the database")
		return models.RentRecord{}, err
	}
	log.Info(spanCtx, "Rent record fetched successfully")
	return rentRecord, nil
}

func (r *rentRecordRepository) GetAllRentRecords(ctx context.Context, userId string,userRole string,rentId string) ([]models.RentRecord, error) {

	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx, "RentRecordRepository.GetAllRentRecords")
	defer span.End()

	rentRecordCollection := r.db.Collection("rent_records")

	log.Info(spanCtx, fmt.Sprintf("Fetching all rent records from the database for user ID: %s with userRole %s", userId, userRole))

	userObjectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		log.Error(spanCtx, "Error converting user ID to ObjectID")
		return nil, err
	}

	var query bson.M

	rentObjectId, err := bson.ObjectIDFromHex(rentId)
	if err != nil {
		log.Error(spanCtx, "Error converting rent ID to ObjectID")
		return nil, err
	}
	query = bson.M{"rent_id": rentObjectId}

	if userRole == string(models.LandLord) {
		query["landlord._id"] = userObjectId
	} else if userRole == string(models.Tenant) {
		query["tenant._id"] = userObjectId
	} else {
		log.Error(spanCtx, "Invalid user role")
		return nil, fmt.Errorf("invalid user role: %s", userRole)
	}


	log.Info(spanCtx, "Querying rent records from the database")

	log.Info(spanCtx, fmt.Sprintf("Query: %v", query))

	cursor, err := rentRecordCollection.Find(spanCtx, query)

	if err != nil {
		log.Error(spanCtx, "Error fetching rent records from the database")
		return nil, err
	}

	defer cursor.Close(spanCtx)

	var rentRecords []models.RentRecord

	if err := cursor.All(spanCtx, &rentRecords); err != nil {
		log.Error(spanCtx, "Error decoding rent records")
		return nil, err
	}

	log.Info(spanCtx, fmt.Sprintf("Found %d rent records", len(rentRecords)))

	return rentRecords, nil
}

func (r *rentRecordRepository) UpdateRentRecord(ctx context.Context, userId string, rentRecordId string, rentRecord models.RentRecord) (models.RentRecord, error) {

	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx, "RentRecordRepository.UpdateRentRecord")
	defer span.End()

	rentRecordCollection := r.db.Collection("rent_records")

	log.Info(spanCtx, "Updating rent record in the database")

	rentRecordObjectId, err := bson.ObjectIDFromHex(rentRecordId)
	if err != nil {
		log.Error(spanCtx, "Error converting rent record ID to ObjectID")
		return models.RentRecord{}, err
	}

	userObjectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		log.Error(spanCtx, "Error converting user ID to ObjectID")
		return models.RentRecord{}, err
	}

	rentId := rentRecord.RentId

	query := bson.M{"_id": rentRecordObjectId, "landlord._id": userObjectId, "rent_id": rentId}

	_, err = rentRecordCollection.UpdateOne(spanCtx, query, bson.M{"$set": rentRecord})
	if err != nil {
		log.Error(spanCtx, "Error updating rent record in the database")
		return models.RentRecord{}, err
	}

	log.Info(spanCtx, "Rent record updated successfully")

	return r.GetRentRecordById(spanCtx, rentRecordId)
}
