package request

type CreateScoreRequest struct {
	StudentID string `json:"student_id" binding:"required"` // ID sinh viên
	SubjectID string `json:"subject_id" binding:"required"` // ID môn học
	Semester  string `json:"semester" binding:"required"`   // Học kỳ

	Attendance float64 `json:"attendance" binding:"required,gte=0,lte=10"` // điểm chuyên cần
	Midterm    float64 `json:"midterm" binding:"required,gte=0,lte=10"`    // điểm giữa kỳ
	Final      float64 `json:"final" binding:"required,gte=0,lte=10"`      // điểm cuối kỳ
}

type CreateScoreByExcelRequest struct {
	StudentCode string  `json:"student_id" binding:"required"`
	SubjectCode string  `json:"subject_code" binding:"required"`
	Semester    string  `json:"semester" binding:"required"`
	Attendance  float64 `json:"attendance" binding:"required"`
	Midterm     float64 `json:"midterm" binding:"required"`
	Final       float64 `json:"final" binding:"required"`
}
