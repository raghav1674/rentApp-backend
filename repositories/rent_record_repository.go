package repositories

import (
	"context"
	"sample-web/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RentRecordRepository interface {
	CreateRentRecord(ctx context.Context, userId string, rentId string, rentRecord models.RentRecord) (models.RentRecord, error)
	GetAllRentRecords(ctx context.Context, userId string, rentId string, userRole string) ([]models.RentRecord, error)
	GetRentRecordById(ctx context.Context, userId string, rentRecordId string) (models.RentRecord, error)
	ApproveRentRecord(ctx context.Context, userId string, rentRecordId string) (models.RentRecord, error)
}

type rentRecordRepository struct {
	db *mongo.Database
}

func NewRentRecordRepository(db *mongo.Database) RentRecordRepository {
	return &rentRecordRepository{
		db: db,
	}
}


func (r *rentRecordRepository) ApproveRentRecord(ctx context.Context, userId string, rentRecordId string) (models.RentRecord, error) {
	panic("unimplemented")
}

func (r *rentRecordRepository) CreateRentRecord(ctx context.Context, userId string, rentId string, rentRecord models.RentRecord) (models.RentRecord, error) {
	panic("unimplemented")
}

func (r *rentRecordRepository) GetAllRentRecords(ctx context.Context, userId string, rentId string, userRole string) ([]models.RentRecord, error) {
	panic("unimplemented")
}

func (r *rentRecordRepository) GetRentRecordById(ctx context.Context, userId string, rentRecordId string) (models.RentRecord, error) {
	panic("unimplemented")
}

