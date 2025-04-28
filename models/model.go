package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRole string

type RentSchedule string

type RentStatus string

type RentRecordStatus string

const (
	LandLord UserRole = "landlord"
	Tenant   UserRole = "tenant"
)

const (
	Monthly   RentSchedule = "monthly"
	BiMonthly RentSchedule = "bimonthly"
)

const (
	RentStatusActive   RentStatus = "active"
	RentStatusInactive RentStatus = "inactive"
)

const (
	RentRecordStatusPending  RentRecordStatus = "pending"
	RentRecordStatusApproved RentRecordStatus = "approved"
	RentRecordStatusRejected RentRecordStatus = "rejected"
)

type User struct {
	Id          bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string        `bson:"name" json:"name"`
	Email       string        `bson:"email" json:"email"`
	Password    string        `bson:"password" json:"password"`
	PhoneNumber string        `bson:"phone_number" json:"phone_number"`
	Roles       []UserRole    `bson:"roles" json:"roles"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
}

type PersonRef struct {
	Id    bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email string        `bson:"email" json:"email"`
}

type RentInfo struct {
	Amount   float64      `bson:"amount" json:"amount"`
	Schedule RentSchedule `bson:"schedule" json:"schedule"`
}

type Rent struct {
	Id        bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	LandLord  PersonRef     `bson:"landlord" json:"landlord"`
	Tenant    PersonRef     `bson:"tenant" json:"tenant"`
	Location  string        `bson:"location" json:"location"`
	Amount    float64       `bson:"amount" json:"amount"`
	Schedule  RentSchedule  `bson:"schedule" json:"schedule"`
	Status    RentStatus    `bson:"status" json:"status"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

type RentRecord struct {
	Id          bson.ObjectID    `bson:"_id,omitempty" json:"_id,omitempty"`
	RentId      bson.ObjectID    `bson:"rent_id" json:"rent_id"`
	Rent        RentInfo         `bson:"rent" json:"rent"`
	Amount      float64          `bson:"amount" json:"amount"`
	SubmittedAt time.Time        `bson:"submitted_at" json:"submitted_at"`
	ApprovedAt  time.Time        `bson:"approved_at" json:"approved_at"`
	Status      RentRecordStatus `bson:"status" json:"status"`
	CreatedAt   time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time        `bson:"updated_at" json:"updated_at"`
}
