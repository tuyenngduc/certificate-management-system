package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID    string             `bson:"studentId" json:"student_id"`
	FullName     string             `bson:"fullName" json:"full_name"`
	Email        string             `bson:"email" json:"email"`
	Ethnicity    string             `bson:"ethnicity" json:"ethnicity"`
	Gender       string             `bson:"gender" json:"gender"`
	FacultyID    primitive.ObjectID `bson:"facultyId" json:"faculty_id"`
	ClassID      primitive.ObjectID `bson:"classId" json:"class_id"`
	Course       string             `bson:"course" json:"course"`
	NationalID   string             `bson:"nationalId" json:"national_id"`
	Address      string             `bson:"address" json:"address"`
	PlaceOfBirth string             `bson:"placeOfBirth" json:"place_of_birth"`
	DateOfBirth  time.Time          `bson:"dateOfBirth" json:"date_of_birth"`
	PhoneNumber  string             `bson:"phoneNumber,omitempty" json:"phone_number,omitempty"`
}
