package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID string             `bson:"studentId" json:"student_id"`
	FullName  string             `bson:"fullName" json:"full_name"`
	Email     string             `bson:"email" json:"email"`
	Faculty   string             `bson:"facultyId" json:"faculty"`
	Class     string             `bson:"classId" json:"class"`
	Course    string             `bson:"course" json:"course"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt" json:"created_at"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updated_at"`
}

type CreateUserRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	FullName  string `json:"full_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Faculty   string `json:"faculty" binding:"required"`
	Class     string `json:"class" binding:"required"`
	Course    string `json:"course" binding:"required,courseyear"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id"`
	StudentID string             `json:"student_id"`
	FullName  string             `json:"full_name"`
	Email     string             `json:"email"`
	Faculty   string             `json:"faculty"`
	Class     string             `json:"class"`
	Course    string             `json:"course"`
	Status    string             `json:"status"`
}

type SearchUserParams struct {
	StudentID string `form:"student_id"`
	FullName  string `form:"full_name"`
	Email     string `form:"email"`
	Class     string `form:"class"`
	Faculty   string `form:"faculty"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
}

type UpdateUserRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	FullName  string `json:"full_name"`
	Email     string `json:"email" binding:"omitempty,email"`
	Faculty   string `json:"faculty"`
	Class     string `json:"class"`
	Course    string `json:"course" binding:"courseyear"`
	Status    string `json:"status"`
}
