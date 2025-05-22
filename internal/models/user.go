package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName     string             `bson:"fullName" json:"full_name" validate:"required,min=3,max=100"`
	Email        string             `bson:"email" json:"email" validate:"required,email"`
	Ethnicity    string             `bson:"ethnicity" json:"ethnicity" validate:"required"`
	Gender       string             `bson:"gender" json:"gender" validate:"required,oneof=male female other"`
	Major        string             `bson:"major" json:"major" validate:"required"`
	Class        string             `bson:"class" json:"class" validate:"required"`
	Course       string             `bson:"course" json:"course" validate:"required"`
	NationalID   string             `bson:"nationalId" json:"national_id" validate:"required,len=12,numeric"`
	Address      string             `bson:"address" json:"address" validate:"required"`
	PlaceOfBirth string             `bson:"placeOfBirth" json:"place_of_birth" validate:"required"`
	DateOfBirth  time.Time          `bson:"dateOfBirth" json:"date_of_birth" validate:"required"`
	PhoneNumber  string             `bson:"phoneNumber,omitempty" json:"phone_number,omitempty" validate:"omitempty,e164"`
}
