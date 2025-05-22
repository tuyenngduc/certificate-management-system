package response

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserLoginResponse struct {
	ID       primitive.ObjectID `json:"id"`
	FullName string             `json:"full_name"`
	Email    string             `json:"email"`
	Role     string             `json:"role"`
	Token    string             `json:"token"`
}
