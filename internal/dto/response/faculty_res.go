package response

import "go.mongodb.org/mongo-driver/bson/primitive"

type FacultyResponse struct {
	ID             primitive.ObjectID `json:"id"`
	Name           string             `json:"name"`
	Code           string             `json:"code"`
	TrainingPeriod string             `json:"training_period"`
}
