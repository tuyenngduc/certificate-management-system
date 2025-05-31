package response

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScoreWithSubjectAndStudentResponse struct {
	ID           primitive.ObjectID `json:"id"`
	StudentName  string             `json:"studentName"`
	SubjectName  string             `json:"subjectName"`
	Semester     string             `json:"semester"`
	Attendance   float64            `json:"attendance"`
	Midterm      float64            `json:"midterm"`
	Final        float64            `json:"final"`
	ProcessScore float64            `json:"processScore"`
	Total        float64            `json:"total"`
	GPAChar      string             `json:"gpaChar"`
	Passed       bool               `json:"passed"`
	Credit       int                `json:"credit"`
}

type CGPAResponse struct {
	CGPA                float64 `json:"cgpa"`
	TotalSubjects       int     `json:"totalSubjects"`
	TotalCredits        int     `json:"totalCredits"`
	TotalFailedSubjects int     `json:"totalFailedSubjects"`
}
