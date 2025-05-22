package request

type CreateUserRequest struct {
	FullName     string `json:"full_name" binding:"required,min=3,max=100"`
	Email        string `json:"email" binding:"required,email"`
	StudentID    string `json:"student_id" binding:"required"`
	Ethnicity    string `json:"ethnicity" binding:"required"`
	Gender       string `json:"gender" binding:"required,oneof=male female other"`
	Major        string `json:"major" binding:"required"`
	Class        string `json:"class" binding:"required"`
	Course       string `json:"course" binding:"required"`
	NationalID   string `json:"national_id" binding:"required,len=12,numeric"`
	Address      string `json:"address" binding:"required"`
	PlaceOfBirth string `json:"place_of_birth" binding:"required"`
	DateOfBirth  string `json:"date_of_birth" binding:"required"` // dd/mm/yyyy
	PhoneNumber  string `json:"phone_number" binding:"omitempty,e164"`
}
