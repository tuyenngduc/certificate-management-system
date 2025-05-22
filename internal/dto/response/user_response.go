package response

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResponse struct {
	ID           primitive.ObjectID `json:"id"`
	FullName     string             `json:"full_name"`
	Email        string             `json:"email"`
	Ethnicity    string             `json:"ethnicity"`
	Gender       string             `json:"gender"`
	Major        string             `json:"major"`
	Class        string             `json:"class"`
	Course       string             `json:"course"`
	NationalID   string             `json:"national_id"`
	Address      string             `json:"address"`
	PlaceOfBirth string             `json:"place_of_birth"`
	DateOfBirth  time.Time          `json:"date_of_birth"`
	PhoneNumber  string             `json:"phone_number,omitempty"`
	Role         string             `json:"role"`
}
