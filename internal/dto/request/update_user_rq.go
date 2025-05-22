package request

type UpdateUserRequest struct {
	FullName     string `json:"full_name" binding:"omitempty,min=3,max=100"`
	Email        string `json:"email" binding:"omitempty,email"`
	StudentID    string `json:"student_id" binding:"omitempty"`
	Ethnicity    string `json:"ethnicity" binding:"omitempty"`
	Gender       string `json:"gender" binding:"omitempty,oneof=male female other"`
	Major        string `json:"major" binding:"omitempty"`
	Class        string `json:"class" binding:"omitempty"`
	Course       string `json:"course" binding:"omitempty"`
	NationalID   string `json:"national_id" binding:"omitempty,len=12,numeric"`
	Address      string `json:"address" binding:"omitempty"`
	PlaceOfBirth string `json:"place_of_birth" binding:"omitempty"`
	DateOfBirth  string `json:"date_of_birth" binding:"omitempty"` // dd/mm/yyyy
	PhoneNumber  string `json:"phone_number" binding:"omitempty,e164"`
}
