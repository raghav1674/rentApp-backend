package models

import "time"

type UserRole string

const (
	LandLord UserRole = "landlord"
	Tenant   UserRole = "tenant"
)

type RentStatus string

const (
	RentStatusActive   RentStatus = "active"
	RentStatusInactive RentStatus = "inactive"
)

type RentRecordStatus string

const (
	RentRecordStatusPending  RentRecordStatus = "pending"
	RentRecordStatusApproved RentRecordStatus = "approved"
	RentRecordStatusRejected RentRecordStatus = "rejected"
)

type User struct {
	Id          string     `bson:"_id,omitempty" json:"id,omitempty"`
	Email       string     `bson:"email" json:"email"`
	Password    string     `bson:"password" json:"password"`
	PhoneNumber string     `bson:"phone_number" json:"phone_number"`
	Roles       []UserRole `bson:"roles" json:"roles"`
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`
}

type PersonRef struct {
	Id    string `bson:"_id,omitempty" json:"id,omitempty"`
	Email string `bson:"email" json:"email"`
}

type Rent struct {
	Id        string     `bson:"_id,omitempty" json:"id,omitempty"`
	LandLord  PersonRef  `bson:"landlord" json:"landlord"`
	Tenant    PersonRef  `bson:"tenant" json:"tenant"`
	Location  string     `bson:"location" json:"location"`
	Amount    float64    `bson:"amount" json:"amount"`
	Schedule  string     `bson:"schedule" json:"schedule"`
	Status    RentStatus `bson:"status" json:"status"`
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
}

type RentRecord struct {
	Id          string           `bson:"_id,omitempty" json:"id,omitempty"`
	RentId      string           `bson:"rent_id" json:"rent_id"`
	Amount      float64          `bson:"amount" json:"amount"`
	SubmittedAt time.Time        `bson:"submitted_at" json:"submitted_at"`
	ApprovedAt  time.Time        `bson:"approved_at" json:"approved_at"`
	Status      RentRecordStatus `bson:"status" json:"status"`
	CreatedAt   time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time        `bson:"updated_at" json:"updated_at"`
}
